package auth

import "time"

const RefreshTokenValidTime = 12 * time.Hour
const AccessTokenValidTime = 60 * time.Minute

const Refresh = "refresh"
const Access = "access"

const CtxKey = "auth"

const ExchangeKey = "exchange_code"
