package td

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrMissingAcctIDs = errors.New("missing account IDs")
)

//go:generate enumer -type UserPrincipalField -trimprefix UserPrincipalField -transform title-lower
type UserPrincipalField byte

const (
	UserPrincipalFieldUnspecified UserPrincipalField = iota
	UserPrincipalFieldStreamerSubscriptionKeys
	UserPrincipalFieldStreamerConnectionInfo
	UserPrincipalFieldPreferences
	UserPrincipalFieldSurrogateIds
)

type Preferences struct {
	ExpressTrading                   bool   `json:"expressTrading"`
	DirectOptionsRouting             bool   `json:"directOptionsRouting"`
	DirectEquityRouting              bool   `json:"directEquityRouting"`
	DefaultEquityOrderLegInstruction string `json:"defaultEquityOrderLegInstruction"`
	DefaultEquityOrderType           string `json:"defaultEquityOrderType"`
	DefaultEquityOrderPriceLinkType  string `json:"defaultEquityOrderPriceLinkType"`
	DefaultEquityOrderDuration       string `json:"defaultEquityOrderDuration"`
	DefaultEquityOrderMarketSession  string `json:"defaultEquityOrderMarketSession"`
	DefaultEquityQuantity            int    `json:"defaultEquityQuantity"`
	MutualFundTaxLotMethod           string `json:"mutualFundTaxLotMethod"`
	OptionTaxLotMethod               string `json:"optionTaxLotMethod"`
	EquityTaxLotMethod               string `json:"equityTaxLotMethod"`
	DefaultAdvancedToolLaunch        string `json:"defaultAdvancedToolLaunch"`
	AuthTokenTimeout                 string `json:"authTokenTimeout"`
}

type StreamerSubscriptionKeys struct {
	Keys []KeyEntry `json:"keys"`
}

type KeyEntry struct {
	Key string `json:"key"`
}

type UserPrincipal struct {
	AuthToken                string                   `json:"authToken"`
	UserID                   string                   `json:"userId"`
	UserCdDomainID           string                   `json:"userCdDomainId"`
	PrimaryAccountID         string                   `json:"primaryAccountId"`
	LastLoginTime            string                   `json:"lastLoginTime"`
	TokenExpirationTime      string                   `json:"tokenExpirationTime"`
	LoginTime                string                   `json:"loginTime"`
	AccessLevel              string                   `json:"accessLevel"`
	StalePassword            bool                     `json:"stalePassword"`
	StreamerInfo             []StreamerInfo           `json:"streamerInfo"`
	ProfessionalStatus       string                   `json:"professionalStatus"`
	Quotes                   QuoteDelays              `json:"quotes"`
	StreamerSubscriptionKeys StreamerSubscriptionKeys `json:"streamerSubscriptionKeys"`
	Accounts                 []UserAccountInfo        `json:"accounts"`
}

type UserAccountInfo struct {
	AccountID         string         `json:"accountId"`
	Description       string         `json:"description"`
	DisplayName       string         `json:"displayName"`
	AccountCdDomainID string         `json:"accountCdDomainId"`
	Company           string         `json:"company"`
	Segment           string         `json:"segment"`
	SurrogateIds      string         `json:"surrogateIds"`
	Preferences       Preferences    `json:"preferences"`
	ACL               string         `json:"acl"`
	Authorizations    Authorizations `json:"authorizations"`
}

type Authorizations struct {
	Apex               bool   `json:"apex"`
	LevelTwoQuotes     bool   `json:"levelTwoQuotes"`
	StockTrading       bool   `json:"stockTrading"`
	MarginTrading      bool   `json:"marginTrading"`
	StreamingNews      bool   `json:"streamingNews"`
	OptionTradingLevel string `json:"optionTradingLevel"`
	StreamerAccess     bool   `json:"streamerAccess"`
	AdvancedMargin     bool   `json:"advancedMargin"`
	ScottradeAccount   bool   `json:"scottradeAccount"`
}

type StreamerInfo struct {
	StreamerSocketURL      string `json:"streamerSocketUrl"`
	SchwabClientCustomerId string `json:"schwabClientCustomerId"`
	SchwabClientCorrelId   string `json:"schwabClientCorrelId"`
	SchwabClientChannel    string `json:"schwabClientChannel"`
	SchwabClientFunctionId string `json:"schwabClientFunctionId"`
}

type QuoteDelays struct {
	IsNyseDelayed   bool `json:"isNyseDelayed"`
	IsNasdaqDelayed bool `json:"isNasdaqDelayed"`
	IsOpraDelayed   bool `json:"isOpraDelayed"`
	IsAmexDelayed   bool `json:"isAmexDelayed"`
	IsCmeDelayed    bool `json:"isCmeDelayed"`
	IsIceDelayed    bool `json:"isIceDelayed"`
	IsForexDelayed  bool `json:"isForexDelayed"`
}

// GetPreferences returns Preferences for a specific account.
// See https://developer.tdameritrade.com/user-principal/apis/get/accounts/%7BaccountId%7D/preferences-0
func (s *HTTPClient) GetPreferences(ctx context.Context, accountID string) (*Preferences, error) {
	preferences := new(Preferences)
	err := s.do(ctx, http.MethodGet, fmt.Sprintf("accounts/%s/preferences", accountID), nil, preferences)
	if err != nil {
		return nil, err
	}

	return preferences, err
}

// GetStreamerSubscriptionKeys returns Subscription Keys for provided accounts or default accounts.
// See https://developer.tdameritrade.com/user-principal/apis/get/userprincipals/streamersubscriptionkeys-0
func (s *HTTPClient) GetStreamerSubscriptionKeys(ctx context.Context, accountIDs ...string) (*StreamerSubscriptionKeys, error) {
	if len(accountIDs) == 0 {
		return nil, ErrMissingAcctIDs
	}

	streamerSubscriptionKeys := new(StreamerSubscriptionKeys)
	err := s.do(
		ctx,
		http.MethodGet,
		fmt.Sprintf("userprincipals/streamersubscriptionkeys?accountIds=%s", strings.Join(accountIDs, ",")),
		nil,
		streamerSubscriptionKeys,
	)

	if err != nil {
		return nil, err
	}

	return streamerSubscriptionKeys, nil
}

// GetUserPrincipals returns User Principal details.
// Valid values for `fields` are "streamerSubscriptionKeys", "streamerConnectionInfo", "preferences" and  "surrogateIds"
// See https://developer.tdameritrade.com/user-principal/apis/get/userprincipals-0
func (s *HTTPClient) GetUserPrincipals(ctx context.Context, fields ...UserPrincipalField) (*UserPrincipal, error) {
	var sb strings.Builder
	if n := len(fields); n > 0 {
		sb.WriteString("?fields=")

		for i, v := range fields {
			sb.WriteString(v.String())
			if i != n-1 {
				sb.WriteRune(',')
			}
		}
	}

	userPrincipal := new(UserPrincipal)
	err := s.do(ctx, http.MethodGet, fmt.Sprintf("userprincipals%s", sb.String()), nil, userPrincipal)
	if err != nil {
		return nil, err
	}

	return userPrincipal, nil
}

// GetUserPreference returns User Preference details.
func (s *HTTPClient) GetUserPreference(ctx context.Context) (*UserPrincipal, error) {
	userPrincipal := new(UserPrincipal)
	err := s.do(ctx, http.MethodGet, "userPreference", nil, userPrincipal)
	if err != nil {
		return nil, err
	}

	return userPrincipal, err
}

// UpdatePreferences updates Preferences for a specific account.
// Please note that the directOptionsRouting and directEquityRouting values cannot be modified via this operation, even though they are in the request body.
// See https://developer.tdameritrade.com/user-principal/apis/put/accounts/%7BaccountId%7D/preferences-0
func (s *HTTPClient) UpdatePreferences(ctx context.Context, accountID string, newPreferences *Preferences) error {
	if accountID == "" {
		return ErrMissingAcctIDs
	}

	if newPreferences == nil {
		return fmt.Errorf("newPreferences is nil")
	}

	err := s.do(ctx, http.MethodPut, fmt.Sprintf("accounts/%s/preferences", accountID), newPreferences, nil)
	if err != nil {
		return err
	}

	return nil
}
