package freeswitch

// # FreeSwitch            PJSIP
// #

var FreeSwitchPJSIPEquivalents = map[string][]string{ //nolint: gochecknoglobals // const is not available for maps in Go
	"username":           {"auth", "username"},
	"realm":              {"auth", "realm"},
	"from-user":          {"endpoint", "from_user"},
	"from-domain":        {"endpoint", "from_domain"},
	"password":           {"auth", "password"},
	"extension":          {"endpoint", "context"},
	"proxy":              {"endpoint", "outbound_proxy"},
	"register-proxy":     {"registration", "outbound_proxy"},
	"expire-seconds":     {"registration", "expiration"},
	"suppress-cng":       {"endpoint", "suppress_cng"},
	"register":           {"registration", "line"},
	"register-transport": {"registration", "transport"},
	"retry-seconds":      {"registration", "retry_interval"},
	"caller-id-in-from":  {"endpoint", "caller_id_in_from"},
	"contact-params":     {"aor", "extra_contact_params"},
	"ping":               {"aor", "qualify_frequency"},
	"ping-max":           {"aor", "ping_max"},
	"ping-min":           {"aor", "ping_min"},
	"rfc-5626":           {"aor", "rfc5626"},
	"reg-id":             {"aor", "reg_id"},
}

// username              auth/username
// realm                 auth/realm
// from-user             endpoint/from_user
// from-domain           endpoint/from_domain
// password              auth/password
// extension             endpoint/context
// proxy                 endpoint/outbound_proxy
// register-proxy        registration/outbound_proxy
// expire-seconds        registration/expiration (probably: endpoint/timers_sess_expire)
// suppress-cng
// register              registration/line (false - no)
// register-transport    transport/protocol
// retry-seconds         registration/retry_interval
// caller-id-in-from     endpoint/symmetric_transport (TODO: Recheck this)
// contact-params        unknown/lines????
// ping                  aor/qualify_frequency
// ping-max
// ping-min
// #
// # ############################################################
// #
// # log-level
// # debug
// # sip-trace
// # context
// # sip-port
// # sip-ip
// # rtp-ip
// # ext-sip-ip
// # dialplan
// # media-option
// # inbound-bypass-media
// # inbound-proxy-media
// # disable-rtp-auto-adjust
// # ignore-183nosdp

// # variables
// # verbose_sdp
// # absolute_codec_string
// # customer_id

// Возвращает максимальное количество полей фрисвич-конфига.
func GetTotalFSKeysNumber() int {
	res := 0
	for range FreeSwitchPJSIPEquivalents {
		res++
	}

	return res
}
