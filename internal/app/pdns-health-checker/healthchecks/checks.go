package healthchecks

import (
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/config"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/eventcheck"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallcheck"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallensurensec3"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/healthchecks/intervallsigningsync"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-health-checker/models"
	"github.com/nameserver-systems/pdns-distribute/pkg/microservice"
)

func StartAllHealthChecks(msobject *microservice.Microservice, configobject *config.ServiceConfiguration,
	actualstate *models.State,
) {
	go eventcheck.StartEventCheckHandler(msobject, configobject, actualstate)
	go intervallcheck.StartPeriodicalCheck(msobject, configobject, actualstate)
	go intervallsigningsync.StartPeridoicalSigningSync(msobject, configobject, actualstate)
	go intervallensurensec3.StartIntervallEnsureNsec3(msobject, configobject, actualstate)
}
