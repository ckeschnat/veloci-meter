package background

import (
	"fmt"
	"log"
	"time"

	"niecke-it.de/veloci-meter/config"
	"niecke-it.de/veloci-meter/icinga"
	"niecke-it.de/veloci-meter/rdb"
	"niecke-it.de/veloci-meter/rules"
)

// normal rules have a warning and a critical, if a rule has an ok value it is special
// in case ok is sset for the rule we only check if there are more mails than defined by ok
// if not an alter level defined by the rule is set

// TODO add info in case the result from icinga is empty; this is beacause of missing check definitions
func CheckForAlerts(config *config.Config, rules *rules.Rules) {
	r := rdb.NewRDB(&config.Redis)
	for {
		// iterate over all rules
		for i, rule := range rules.Rules {
			actCount := r.CountMail(rule.Pattern)
			//r.RemoveAllAlert(rule.Pattern)
			if rule.Ok != 0 {
				if actCount < rule.Ok {
					if rule.Alert == "critical" {
						//r.StoreAlert(config.AlertInterval, rule.Pattern, "critical")
						icinga.SendResults(config, rule.Name, rule.Pattern, 2)
					} else {
						//r.StoreAlert(config.AlertInterval, rule.Pattern, "warning")
						icinga.SendResults(config, rule.Name, rule.Pattern, 1)
					}
				} else {
					//r.StoreAlert(config.AlertInterval, rule.Pattern, "ok")
					icinga.SendResults(config, rule.Name, rule.Pattern, 0)
				}
			} else {
				// remove all alerts if there are any
				// the alert will be set again in each iteration
				if actCount > rule.Critical {
					//r.StoreAlert(config.AlertInterval, rule.Pattern, "critical")
					icinga.SendResults(config, rule.Name, rule.Pattern, 2)
					log.Println("Rule " + fmt.Sprint(i) + " is CRITICAL: " + rules.Rules[i].ToString())
				} else if actCount > rule.Warning {
					//r.StoreAlert(config.AlertInterval, rule.Pattern, "warning")
					icinga.SendResults(config, rule.Name, rule.Pattern, 1)
					log.Println("Rule " + fmt.Sprint(i) + " is WARNING: " + rules.Rules[i].ToString())
				} else {
					//r.StoreAlert(config.AlertInterval, rule.Pattern, "ok")
					icinga.SendResults(config, rule.Name, "Everythin is fine.", 0)
					log.Println("Rule " + fmt.Sprint(i) + " is OK: " + rules.Rules[i].ToString())
				}
			}
		}

		log.Printf("Sleep for %ds before checking again the threshholds", config.CheckInterval)
		time.Sleep(time.Duration(config.CheckInterval) * time.Second)
	}
}
