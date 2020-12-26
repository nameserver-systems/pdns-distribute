package healthchecks

import (
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/eventcheck"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallcheck"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallensurensec3"

	// nolint:lll
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallsigningsync"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

func StartAllHealthChecks(msobject *microservice.Microservice, configobject *config.ServiceConfiguration,
	actualstate *models.State) {
	go eventcheck.StartEventCheckHandler(msobject, configobject, actualstate)
	go intervallcheck.StartPeriodicalCheck(msobject, configobject, actualstate)
	go intervallsigningsync.StartPeridoicalSigningSync(msobject, configobject, actualstate)
	go intervallensurensec3.StartIntervallEnsureNsec3(msobject, configobject, actualstate)
}
