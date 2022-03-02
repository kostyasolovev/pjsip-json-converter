package pjsip

type Transport struct {
	Name                     []string `json:"name,omitempty" sip:"general,name"`
	Type                     string   `json:"type,omitempty" sip:",hide"`
	Bind                     string   `json:"bind,omitempty" sip:"udpbindaddr,"`
	Protocol                 string   `json:"protocol,omitempty" sip:",omit"`
	LocalNet                 string   `json:"localnet,omitempty" sip:",omit"`
	ExternalMediaAddress     string   `json:"externalmediaaddress,omitempty" sip:",omit"`
	ExternalSignalingAddress string   `json:"externalsignalingaddress,omitempty" sip:",omit"`
	ExternalSignalingPort    string   `json:",omitempty" sip:",omit"`
	CertFile                 string   `json:"certfile,omitempty" sip:",omit"`
	PrivKeyFile              string   `json:"privkeyfile,omitempty" sip:",omit"`
	Cipher                   string   `json:"cipher,omitempty" sip:",omit"`
	Method                   string   `json:"method,omitempty" sip:",omit"`
	AsyncOperations          string   `json:",omitempty" sip:",omit"`
	CaListFile               string   `json:",omitempty" sip:",omit"`
	CaListPath               string   `json:",omitempty" sip:",omit"`
	Domain                   string   `json:",omitempty" sip:",omit"`
	Password                 string   `json:",omitempty" pjsip:"omit" sip:",omit"`
	RequireClientCert        string   `json:",omitempty" sip:",omit"`
	VerifyClient             string   `json:",omitempty" sip:",omit"`
	VerifyServer             string   `json:",omitempty" sip:",omit"`
	Tos                      string   `json:",omitempty" sip:",omit"`
	Cos                      string   `json:",omitempty" sip:",omit"`
	WebsocketWriteTimeout    string   `json:",omitempty" sip:",omit"`
	AllowReload              string   `json:",omitempty" sip:",omit"`
	SymmetricTransport       string   `json:",omitempty" sip:",omit"`
}

type Trunk struct {
	Name         string          `json:",omitempty" fswitch:"gateway" sip:",hide"`
	Endpoint     Endpoint        `json:"endpoint,omitempty" fswitch:"omitempty" sip:",hide"`
	Aor          Aor             `json:"aor,omitempty" fswitch:"omitempty" sip:",hide"`
	Auth         Auth            `json:"auth,omitempty" fswitch:"omitempty" sip:",hide"`
	Registration Registration    `json:"registration,omitempty" fswitch:"omitempty" sip:",omit"`
	Identify     Identify        `json:"identify,omitempty" sip:",omit"`
	Acl          Acl             `json:",omitempty" sip:",hide"` // nolint: revive, stylecheck
	Unknowns     []UnknownStruct `json:",omitempty" sip:",omit"`
}

type Aor struct {
	Name                []string `json:"name,omitempty" sip:",hide"`
	Type                string   `json:"type,omitempty" sip:",hide"`
	MaxContacts         string   `json:"maxcontacts,omitempty" sip:",omit"`
	Contact             []string `json:"contact,omitempty" sip:",omit"`
	MaximumExpiration   string   `json:",omitempty" sip:",omit"`
	MinimumExpiration   string   `json:",omitempty" sip:",omit"`
	RemoveExisting      string   `json:",omitempty" sip:",omit"`
	QualifyFrequency    string   `json:",omitempty" fswitch:"ping" sip:"qualifyfreq,"`
	QualifyTimeout      string   `json:",omitempty" sip:"qualify,"`
	AuthenticateQualify string   `json:",omitempty" sip:",omit"`
	OutboundProxy       string   `json:",omitempty" sip:",omit"`
	Mailboxes           string   `json:",omitempty" sip:"mailbox,"`
	DefaultExpiration   string   `json:",omitempty" sip:",omit"`
	PingMin             string   `json:",omitempty" pjsip:"omit" fswitch:"ping-min" sip:",omit"`
	PingMax             string   `json:",omitempty" pjsip:"omit" fswitch:"ping-max" sip:",omit"`
	ExtraContactParams  string   `json:",omitempty" pjsip:"omit" fswitch:"contact-params" sip:",omit"`
	RegId               string   `json:",omitempty" pjsip:"omit" fswitch:"reg-id" sip:",omit"` // nolint: revive, stylecheck, nolintlint
	TestFailField       string   `fswitch:",failtestcase" sip:",omit"`
	Rfc5626             string   `pjsip:"omit" fswitch:"rfc-5626" sip:",omit"`
}

type Identify struct {
	Name        []string `json:"name,omitempty" sip:",omit"`
	Type        string   `json:"type,omitempty" sip:",omit"`
	Match       []string `json:"match,omitempty" sip:",omit"`
	MatchHeader []string `json:",omitempty" sip:",omit"`
	SrvLookups  string   `json:",omitempty" sip:",omit"`
	Endpoint    string   `json:",omitempty" sip:",omit"`
}
type Endpoint struct {
	FromUser                        string   `json:",omitempty" fswitch:"from-user" sip:",omit"`
	Name                            []string `json:"name,omitempty" sip:",hide"`
	Type                            string   `json:"type,omitempty" sip:",hide"`
	Allow                           []string `json:"allow,omitempty" sip:"allow,"`
	AllowOverlap                    string   `json:",omitempty" sip:",omit"`
	Aors                            string   `json:"aors,omitempty" sip:",omit"`
	Auth                            string   `json:"auth,omitempty" sip:",omit"`
	OutboundAuth                    string   `json:",omitempty" sip:",omit"`
	Callerid                        string   `json:",omitempty" sip:"callerid,"`
	CalleridPrivacy                 string   `json:",omitempty" sip:",omit"`
	CalleridTag                     string   `json:",omitempty" sip:",omit"`
	Context                         string   `json:"context,omitempty" fswitch:"extension" sip:"context,"`
	DirectMediaGlareMitigation      string   `json:",omitempty" sip:",omit"`
	DirectMediaMethod               string   `json:",omitempty" sip:",omit"`
	TrustConnectedLine              string   `json:",omitempty" sip:",omit"`
	SendConnectedLine               string   `json:",omitempty" sip:",omit"`
	ConnectedLineMethod             string   `json:",omitempty" sip:",omit"`
	DirectMedia                     string   `json:"directmedia,omitempty" sip:"directmedia,"`
	DisableDirectMediaOnNat         string   `json:",omitempty" sip:",omit"`
	Disallow                        []string `json:"disallow,omitempty" sip:"disallow,"`
	DtmfMode                        string   `json:",omitempty" sip:"dtfmmode,"`
	MediaAddress                    string   `json:",omitempty" sip:",omit"`
	BindRtpToMediaAddress           string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive,nolintlint
	ForceRport                      string   `json:"forcerport,omitempty" sip:",hide"`
	IceSupport                      string   `json:"icesupport,omitempty" sip:",omit"`
	IdentifyBy                      string   `json:",omitempty" sip:",omit"`
	RedirectMethod                  string   `json:",omitempty" sip:",omit"`
	Mailboxes                       string   `json:"mailboxes,omitempty" sip:",omit"`
	VoicemailExtension              string   `json:",omitempty" sip:"vmexten,"`
	MwiSubscribeReplacesUnsolicited string   `json:",omitempty" sip:",omit"`
	MohSuggest                      string   `json:",omitempty" sip:",omit"`
	MohPassthrough                  string   `json:",omitempty" sip:",omit"`
	OutboundProxy                   string   `json:",omitempty" fswitch:"proxy" sip:",omit"`
	RewriteContact                  string   `json:"rewritecontact,omitempty" sip:",hide"`
	RtpSymmetric                    string   `json:"rtpsymmetric,omitempty" sip:",hide"` // nolint: stylecheck,revive,nolintlint
	SendDiversion                   string   `json:",omitempty" sip:",omit"`
	SendPai                         string   `json:",omitempty" sip:",omit"`
	SendRpid                        string   `json:",omitempty" sip:",omit"`
	RpidImmediate                   string   `json:",omitempty" sip:",omit"`
	TimersMinSe                     string   `json:",omitempty" sip:",omit"`
	Timers                          string   `json:",omitempty" sip:",omit"`
	TimersSessExpires               string   `json:",omitempty" sip:",omit"`
	Transport                       []string `json:"transport,omitempty" sip:",omit"`
	TrustIdInbound                  string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive
	TrustIdOutbound                 string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive
	UsePtime                        string   `json:",omitempty" sip:",omit"`
	UseAvpf                         string   `json:",omitempty" sip:",omit"`
	MediaEncryption                 []string `json:"mediaencryption,omitempty" sip:",omit"`
	MediaEncryptionOptimistic       string   `json:",omitempty" sip:",omit"`
	G726NonStandart                 string   `json:",omitempty" sip:",omit"`
	InbandProgress                  string   `json:",omitempty" sip:",omit"`
	CallGroup                       string   `json:",omitempty" sip:"callgroup,"`
	PickupGroup                     string   `json:",omitempty" sip:"pickupgroup,"`
	NamedCallGroup                  string   `json:",omitempty" sip:"namedcallgroup,"`
	NamedPickupGroup                string   `json:",omitempty" sip:"namedpickupgroup,"`
	DeviceStateBusyAt               string   `json:",omitempty" sip:",omit"`
	AggregateMwi                    string   `json:"aggregatemwi,omitempty" sip:",omit"`
	DeviceStateBusy                 string   `json:"devicestatebusy,omitempty" sip:",omit"`
	StirShaken                      string   `json:"stirshaken,omitempty" sip:",omit"`
	T38Udptl                        string   `json:",omitempty" sip:",omit"`
	T38UdptlEc                      string   `json:",omitempty" sip:",omit"`
	T38UdptlMaxdatagram             string   `json:",omitempty" sip:",omit"`
	FaxDetect                       string   `json:",omitempty" sip:",omit"`
	FaxDetectTimeout                string   `json:",omitempty" sip:",omit"`
	T38UdptlNat                     string   `json:",omitempty" sip:",omit"`
	ToneZone                        string   `json:",omitempty" sip:",omit"`
	Language                        string   `json:",omitempty" sip:"language,"`
	OneTouchRecording               string   `json:",omitempty" sip:",omit"`
	RecordOnFeature                 string   `json:",omitempty" sip:",omit"`
	RecordOffFeature                string   `json:",omitempty" sip:",omit"`
	RtpEngine                       string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive,nolintlint
	AllowTransfer                   string   `json:",omitempty" sip:",omit"`
	SdpOwner                        string   `json:",omitempty" sip:",omit"`
	SdpSession                      string   `json:",omitempty" sip:",omit"`
	TosAudio                        string   `json:",omitempty" sip:",omit"`
	TosVideo                        string   `json:",omitempty" sip:",omit"`
	CosAudio                        string   `json:",omitempty" sip:",omit"`
	CosVideo                        string   `json:",omitempty" sip:",omit"`
	AllowSubscribe                  string   `json:"allowsubscribe,omitempty" sip:",omit"`
	SubMinExpiry                    string   `json:"subminexpiry,omitempty" sip:",omit"`
	MwiFromUser                     string   `json:"mwifromuser,omitempty" sip:",omit"`
	FromDomain                      string   `json:",omitempty" fswitch:"from-domain" sip:",omit"`
	DtlsVerify                      string   `json:",omitempty" sip:",omit"`
	DtlsRekey                       string   `json:",omitempty" sip:",omit"`
	DtlsAutoGenerateCert            string   `json:",omitempty" sip:",omit"`
	DtlsCertFile                    string   `json:",omitempty" sip:",omit"`
	DtlsPrivateKey                  string   `json:",omitempty" sip:",omit"`
	DtlsCipher                      string   `json:",omitempty" sip:",omit"`
	DtlsCaFile                      string   `json:",omitempty" sip:",omit"`
	DtlsCaPath                      string   `json:",omitempty" sip:",omit"`
	DtlsSetup                       string   `json:",omitempty" sip:",omit"`
	DtlsFingerprint                 string   `json:",omitempty" sip:",omit"`
	SrtpTag_32                      string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive
	SetVar                          string   `json:",omitempty" sip:",omit"`
	RtpKeepalive                    string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive,nolintlint
	RtpTimeout                      string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive,nolintlint
	RtpTimeoutHold                  string   `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive,nolintlint
	ContactUser                     string   `json:",omitempty" sip:",omit"`
	IncomingCallOfferPref           string   `json:",omitempty" sip:",omit"`
	OutgoingCallOfferPref           string   `json:",omitempty" sip:",omit"`
	SuppressCng                     string   `json:",omitempty" pjsip:"omit" fswitch:"suppress-cng" sip:",omit"`
	CallerIdInFrom                  string   `json:",omitempty" pjsip:"omit" fswitch:"caller-id-in-from" sip:",omit"` // nolint: revive, stylecheck, nolintlint, lll
	// 100rel                         string   // 100rel - единственное поле, начинающееся с цифр.
	// Поэтому с его имплементацией в го возникают трудности. Можно назвать его подходящим для Го образом,
	// например ХХХХ100rel, но тогда нужно разрабатывать механизм обработки этого кейса и его перевода в
	// в корректное поле pjsip 100rel. Учитывая, что это единственное проблемное поле, притом не самое
	// часто использующееся, решение проблемы пока отложено
}

type Auth struct {
	Name          []string `json:"name,omitempty" sip:",omit"`
	Type          string   `json:"type,omitempty" sip:",omit"`
	NonceLifetime string   `json:",omitempty" sip:",omit"`
	Md5Cred       string   `json:",omitempty" sip:",omit"`
	AuthType      string   `json:",omitempty" sip:",omit"`
	Password      string   `json:",omitempty" fswitch:"password" sip:"secret,"`
	Username      string   `json:",omitempty" fswitch:"username" sip:"username,"`
	Realm         string   `json:",omitempty" fswitch:"realm" sip:",omit"`
}

type PJSIP struct {
	Transports   []Transport `json:",omitempty" sip:",hide"`
	Phoneprovs   []Phoneprov `json:",omitempty" sip:",omit"`
	System       `json:",omitempty" sip:",omit"`
	Global       `json:",omitempty" sip:",omit"`
	ResourceList `json:",omitempty" sip:",omit"`
	Domains      []Domain `json:",omitempty" sip:",omit"`
	// Outbounds []OutboundPublish //DELETED
	Trunks []Trunk `json:",omitempty" sip:",hide"`
	// Unknowns []UnknownStruct `json:",omitempty"` // REPLACED
}

type Registration struct {
	Name                   []string `json:"name,omitempty" sip:",omit"`
	Type                   string   `json:"type,omitempty" sip:",omit"`
	Transport              []string `json:"transport,omitempty" fswitch:"register-transport" sip:",omit"`
	Endpoint               []string `json:",omitempty" sip:",omit"`
	OutboundAuth           string   `json:"outboundauth,omitempty" sip:",omit"`
	ServerUri              string   `json:"serveruri,omitempty" sip:",omit"` // nolint: stylecheck,revive
	AuthRejectionPermanent string   `json:",omitempty" sip:",omit"`
	ClientUri              string   `json:"clienturi,omitempty" sip:",omit"` // nolint: stylecheck,revive
	ContactUser            string   `json:"contactuser,omitempty" sip:",omit"`
	RetryInterval          string   `json:"retryinterval,omitempty" fswitch:"retry-seconds" sip:",omit"`
	ForbiddenRetryInterval string   `json:"forbiddenretryinterval,omitempty" sip:",omit"`
	FatalRetryInterval     string   `json:",omitempty" sip:",omit"`
	Expiration             string   `json:"expiration,omitempty" fswitch:"expire-seconds" sip:",omit"`
	MaxRetries             string   `json:",omitempty" sip:",omit"`
	OutboundProxy          string   `json:",omitempty" fswitch:"register-proxy" sip:",omit"`
	Line                   string   `json:"line,omitempty" fswitch:"register" sip:",omit"`
}

type Acl struct { // nolint: stylecheck,revive
	Name          []string `json:",omitempty" sip:",omit"`
	Type          string   `json:",omitempty" sip:",omit"`
	ContactAcl    []string `json:",omitempty" sip:",omit"` // nolint: stylecheck,revive
	Deny          []string `json:",omitempty" sip:"deny,"`
	Permit        []string `json:",omitempty" sip:"permit,"`
	ContactDeny   []string `json:",omitempty" sip:",omit"`
	ContactPermit []string `json:",omitempty" sip:",omit"`
}

type ResourceList struct {
	Name                      []string `json:",omitempty"`
	Type                      string   `json:",omitempty"`
	ListItem                  []string `json:",omitempty"`
	Event                     []string `json:",omitempty"`
	FullState                 string   `json:",omitempty"`
	NotificationBatchInterval string   `json:",omitempty"`
}

type Phoneprov struct {
	Name     []string `json:",omitempty"`
	Type     string   `json:",omitempty"`
	Endpoint []string `json:",omitempty"`
	PROFILE  string   `json:",omitempty"`
	MAC      string   `json:",omitempty"`
	SERVER   string   `json:",omitempty"`
	TIMEZONE string   `json:",omitempty"`
	MYVAR    string   `json:",omitempty"`
	LABEL    string   `json:",omitempty"`
	OTHERVAR string   `json:",omitempty"`
}

type Domain struct {
	Name   string `json:",omitempty"`
	Type   string `json:",omitempty"`
	Domain string `json:",omitempty"`
}

type System struct {
	Name                     string `json:",omitempty"`
	Type                     string `json:",omitempty"`
	TimerT1                  string `json:",omitempty"`
	TimerB                   string `json:",omitempty"`
	CompactHeaders           string `json:",omitempty"`
	ThreadpoolInitialSize    string `json:",omitempty"`
	ThreadpoolAutoIncrement  string `json:",omitempty"`
	ThreadpoolIdleTimeout    string `json:",omitempty"`
	ThreadpoolMaxSize        string `json:",omitempty"`
	DisableTcpSwitch         string `json:",omitempty"` // nolint: stylecheck,revive
	FollowEarlyMediaFork     string `json:",omitempty"`
	AcceptMultipleSdpAnswers string `json:",omitempty"`
	DisableRport             string `json:",omitempty"`
}

type Global struct {
	MaxForwards                           string `json:",omitempty"`
	Name                                  string `json:",omitempty"`
	Type                                  string `json:",omitempty"`
	UserAgent                             string `json:",omitempty"`
	DefaultOutboundEndpoint               string `json:",omitempty"`
	Debug                                 string `json:",omitempty"`
	KeepAliveInterval                     string `json:",omitempty"`
	ContactExpirationCheckInterval        string `json:",omitempty"`
	DisableMultiDomain                    string `json:",omitempty"`
	EndpointIdentifierOrder               string `json:",omitempty"`
	MaxInitialQualifyTime                 string `json:",omitempty"`
	Regcontext                            string `json:",omitempty"`
	DefaultVoicemailExtension             string `json:",omitempty"`
	UnidentifiedRequestCount              string `json:",omitempty"`
	UnidentifiedRequestPeriod             string `json:",omitempty"`
	UnidentifiedRequestPruneInterval      string `json:",omitempty"`
	DefaultFromUser                       string `json:",omitempty"`
	DefaultRealm                          string `json:",omitempty"`
	MwiTpsQueueHigh                       string `json:",omitempty"`
	MwiTpsQueueLow                        string `json:",omitempty"`
	MwiDisableInitialUnsolicited          string `json:",omitempty"`
	IgnoreUriUserOptions                  string `json:",omitempty"` // nolint: stylecheck,revive
	SendContactStatusOnUpdateRegistration string `json:",omitempty"`
	TaskprocessorOverloadTrigger          string `json:",omitempty"`
	Norefersub                            string `json:",omitempty"`
}

type OutboundPublish struct {
	Name            []string `json:",omitempty"`
	Type            string   `json:",omitempty"`
	Expiration      string   `json:",omitempty"`
	OutboundAuth    string   `json:",omitempty"`
	OutboundProxy   string   `json:",omitempty"`
	ServerUri       string   `json:",omitempty"` // nolint: stylecheck,revive
	FromUri         string   `json:",omitempty"` // nolint: stylecheck,revive
	ToUri           string   `json:",omitempty"` // nolint: stylecheck,revive
	Event           string   `json:",omitempty"`
	MaxAuthAttempts string   `json:",omitempty"`
	Transport       []string `json:",omitempty"`
	MultiUser       string   `json:",omitempty"`
}

type UnknownStruct struct {
	Name  []string `json:",omitempty"`
	Type  string   `json:",omitempty"`
	LINES []string `json:",omitempty"`
}
