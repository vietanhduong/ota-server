package profile

import (
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/metadata"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/storage_object"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/user"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/logger"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/mysql/models"
	"github.com/vietanhduong/ota-server/pkg/notifications/telegram"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"gopkg.in/errgo.v2/errors"
	"net/http"
)

const AppIcon = "app_icon"

type StorageService interface {
	GetObjectById(objectId int) (*storage_object.File, error)
	GetObjectsByKeys(objectKeys []string) (map[string]*storage_object.File, error)
}

type MetadataService interface {
	CreateMetadata(profileId int, metadata map[string]string) ([]*metadata.Metadata, error)
	GetMetadata(profileId int) ([]*metadata.Metadata, error)
	GetMetadataByListProfileId(profileIds []uint) (map[uint][]*metadata.Metadata, error)
}

type UserRepository interface {
	FindByIds(userIds []int, active bool) (map[int]*models.User, error)
	FindById(userId uint) (*models.User, error)
}

type service struct {
	repo        *repository
	telegramSvc *telegram.Telegram
	storageSvc  StorageService
	metadataSvc MetadataService
	userRepo    UserRepository
}

func NewService(db *mysql.DB) *service {
	var _telegram *telegram.Telegram
	telegramToken := env.GetEnvAsStringOrFallback("TELEGRAM_BOT_TOKEN", "")
	telegramGroupId := env.GetEnvAsStringOrFallback("TELEGRAM_GROUP_ID", "")
	if telegramToken == "" || telegramGroupId == "" {
		logger.Logger.Warnf("not found telegram bot token or telegram group id in environment variables => STOP initialize telegram")
	} else {
		_telegram = telegram.InitializeTelegram(telegramToken, telegramGroupId)
	}

	return &service{
		repo:        NewRepository(db),
		storageSvc:  storage_object.NewService(db),
		metadataSvc: metadata.NewService(db),
		telegramSvc: _telegram,
		userRepo:    user.NewRepository(db),
	}
}

func (s *service) GetProfiles() ([]*ResponseProfile, error) {
	profiles, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	// prepare profile ids
	var profileIds []uint
	var userIds []int
	for _, p := range profiles {
		profileIds = append(profileIds, p.ID)
		userIds = append(userIds, int(p.UserID))
	}

	// fetch user
	users, err := s.userRepo.FindByIds(userIds, true)
	if err != nil {
		return nil, err
	}

	// fetch metadata
	mm, err := s.metadataSvc.GetMetadataByListProfileId(profileIds)
	if err != nil {
		return nil, err
	}

	var metadataList []*metadata.Metadata
	for _, ml := range mm {
		metadataList = append(metadataList, ml...)
	}

	appIcons, err := s.GetAppIconInMetadata(metadataList)
	if err != nil {
		return nil, err
	}

	// convert to response object
	var result []*ResponseProfile
	for _, p := range profiles {
		if u, ok := users[int(p.UserID)]; ok {
			p.User = u
		}

		profile := ToResponseProfile(p)
		if m, ok := mm[profile.ProfileId]; ok {
			profile.Metadata = ConvertMetadataListToMap(m)

			// hook download icon url
			if v, ok := profile.Metadata[AppIcon]; ok {
				if url, found := appIcons[v]; found {
					profile.Metadata[AppIcon] = url
				}
			}
		}

		result = append(result, profile)
	}

	return result, nil
}

func (s *service) GetProfile(profileId int) (*ResponseProfile, error) {
	model, err := s.repo.FindById(uint(profileId))
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, cerrors.NewCError(http.StatusNotFound, "profile does not exist")
	}

	// find object by object id
	object, err := s.storageSvc.GetObjectById(int(model.StorageObjectID))
	if err != nil {
		return nil, err
	}

	// convert to object model
	// this object is used in generating a manifest
	model.StorageObject = &models.StorageObject{
		Key:  object.Key,
		Name: object.Filename,
	}

	createdBy, err := s.userRepo.FindById(model.UserID)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if createdBy != nil {
		model.User = &models.User{
			Email:       createdBy.Email,
			DisplayName: createdBy.DisplayName,
		}
	}

	profile := ToResponseProfile(model)

	m, err := s.metadataSvc.GetMetadata(profileId)
	if err != nil {
		return nil, err
	}

	profile.Metadata = ConvertMetadataListToMap(m)

	appIcon, err := s.GetAppIconInMetadata(m)
	if err != nil {
		return nil, err
	}

	if len(appIcon) != 0 {
		profile.Metadata[AppIcon] = appIcon[profile.Metadata[AppIcon]]
	}

	return profile, nil
}

func (s *service) CreateProfile(reqProfile *RequestProfile) (*ResponseProfile, error) {
	// TODO: update validate before insert to database
	// validate storage object
	_, err := s.storageSvc.GetObjectById(reqProfile.StorageObjectID)
	if err != nil {
		return nil, err
	}
	// insert to database
	profileModel, err := s.repo.Insert(reqProfile)
	if err != nil {
		return nil, err
	}

	profile := ToResponseProfile(profileModel)

	if len(reqProfile.Metadata) > 0 {
		m, err := s.metadataSvc.CreateMetadata(int(profileModel.ID), reqProfile.Metadata)
		if err != nil {
			return nil, err
		}
		profile.Metadata = ConvertMetadataListToMap(m)
	}

	// notify to Telegram
	// to avoid the main thread, please separate
	// another thread to send notifications to telegram
	go func() {
		// stop if telegram service is not initialized
		if s.telegramSvc == nil {
			return
		}
		msg := createNotificationMessage(profile)
		if err := s.telegramSvc.SendMessage(msg); err != nil {
			logger.Logger.Errorf("send message to telegram failed with err: %v", err)
		}
	}()

	return profile, nil
}

func createNotificationMessage(profile *ResponseProfile) string {
	// new line character
	const newLine = "\n"

	host := env.GetEnvAsStringOrFallback("HOST", "https://ota.anhdv.dev")
	title := fmt.Sprintf("*[iOS]* Just got a new *build (#%d)* uploaded to OTA server [%s](%s)", profile.ProfileId, host, host)
	info := fmt.Sprintf("*Information*%s---%s*App name:*` %s` %s*Version:*` %s` %s*Build:*` %d`", newLine, newLine, profile.AppName, newLine, profile.Version, newLine, profile.Build)
	// stop send git information if repo is not appeared in metadata
	repo, found := profile.Metadata["repo"]
	if !found {
		return fmt.Sprintf("%s%s%s", title, newLine, info)
	}

	// send git commit info
	var git string
	if commit, found := profile.Metadata["commit"]; found && len(commit) > 6 {
		git = fmt.Sprintf("*Commit:* [%s](%s/commit/%s)", commit[:6], repo, commit)
	}

	// send pull request info
	if prNumber, found := profile.Metadata["pr_number"]; found {
		git = fmt.Sprintf("*PR:* [#%s](%s/pull/%s)", prNumber, repo, prNumber)
	}

	return fmt.Sprintf("%s%s%s%s%s", title, newLine, info, newLine, git)
}

// GetAppIconInMetadata get app icon in input list metadata
// and return a map contains objectKey - download url
func (s *service) GetAppIconInMetadata(metadataList []*metadata.Metadata) (map[string]string, error) {
	result := make(map[string]string)
	var objectKeys []string

	// get all object key in metadata
	for _, m := range metadataList {
		if m.Key == AppIcon {
			objectKeys = append(objectKeys, m.Value)
		}
	}
	// if there are no object key, just end this function
	if len(objectKeys) == 0 {
		return result, nil
	}

	// get all objects
	// this function return a map contain objectKey - storage_object.File
	objects, err := s.storageSvc.GetObjectsByKeys(objectKeys)
	if err != nil {
		return nil, err
	}

	// generate download url
	for key, obj := range objects {
		result[key] = fmt.Sprintf("%s/api/v1/storages/%s/download/%s", Host, key, obj.Filename)
	}

	return result, nil
}
