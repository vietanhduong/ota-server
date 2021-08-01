import React from 'react';
import {
  Avatar,
  Box,
  Button,
  Link,
  List,
  ListItem,
  ListItemAvatar,
  ListItemSecondaryAction,
  ListItemText,
  Paper,
  Typography
} from '@material-ui/core';
import {profileService} from 'services/profile';
import {getDownloadUrl, getExchangeCode} from 'utils/common';
import {profileAction} from 'actions/profile';

function Home() {
  const [data, setData] = React.useState([]);

  const fetchData = React.useCallback(() => {
    profileService
      .getProfiles()
      .then((res) => {
        setData(res || []);
      })
      .catch((e) => console.log(e));
  }, []);

  const renderGitInfo = (metadata) => {
    metadata = metadata || {};
    // just ignore if metadata does not contains `repo`, `commit` or `pr_number` key
    if (!('repo' in metadata) || (!('commit' in metadata) && !('pr_number' in metadata))) return '';
    // if `pr_number` is null => render commit hash
    if (!metadata.pr_number) {
      return (
        <span className={'metadata-attribute'}>
          <b className={'mr5'}>commit:</b>
          <Link href={`${metadata.repo}/commit/${metadata.commit}`} target='_blank'>
            {metadata.commit.substring(0, 6)}
          </Link>
        </span>
      );
    }
    // else render pr number
    return (
      <span className={'metadata-attribute'}>
        <b className={'mr5'}>pr:</b>
        <Link href={`${metadata.repo}/pull/${metadata.pr_number}`} target='_blank'>
          #{metadata.pr_number}
        </Link>
      </span>
    );
  };

  const renderNoData = (d) => {
    if (d.length > 0) return '';
    return (
      <ListItem>
        <ListItemText
          primary={'No data available'}
          style={{textAlign: 'center', fontStyle: 'italic', color: 'gray'}}
        />
      </ListItem>
    );
  };

  const exchangeCode = getExchangeCode();

  React.useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <Box style={{}}>
      <Box style={{maxWidth: 680, margin: '0 auto'}}>
        <List disablePadding component={Paper} style={{marginBottom: 10}} variant='outlined'>
          <ListItem>
            <ListItemText
              primary={
                <Box display='flex' justifyContent='space-between'>
                  <Box display='flex' alignItems='center'>
                    <img alt='' style={{height: 28}} src={'apple-icon-57x57.png'}/>
                    <Typography variant='h5' style={{paddingTop: 3}}>
                      Over-The-Air Server
                    </Typography>
                  </Box>
                  <Button color='secondary' onClick={profileAction.logout}>
                    <Typography variant='body2' style={{fontWeight: "bold"}}>Logout</Typography>
                  </Button>
                </Box>
              }
            />
          </ListItem>
        </List>
        <List disablePadding component={Paper} variant='outlined'>
          {renderNoData(data)}
          {data.map((item, index) => (
            <ListItem key={item.profile_id} divider={index + 1 !== data.length} style={{paddingRight: 80}}>
              <ListItemAvatar>
                <Avatar alt="Travis Howard" src="ios-app-icon.png"/>
              </ListItemAvatar>
              <ListItemText
                primary={
                  <span style={{fontSize: ".9rem", fontWeight: "bold"}}>
                    #{item.profile_id}: {item.app_name}
                  </span>
                }
                secondary={
                  <>
                    <span className={'metadata-attribute'}>
                      <b className={'mr5'}>version:</b> {item.version}
                    </span>
                    <span className={'metadata-attribute'}>
                      <b className={'mr5'}>build:</b> {item.build}
                    </span>
                    {renderGitInfo(item.metadata)}
                  </>
                }
              />
              <ListItemSecondaryAction>
                <Link href={getDownloadUrl(item, exchangeCode)} target='_blank'>
                  <Button disableElevation variant='contained' color='primary' style={{borderRadius: 18}}>
                    <Typography variant='body2'>GET</Typography>
                  </Button>
                </Link>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
      </Box>
    </Box>
  );
}

export default Home;
