package main

import (
	"fmt"
	"time"

	"github.com/crowdsecurity/crowdsec/pkg/sqlite"
	"github.com/crowdsecurity/crowdsec/pkg/types"
)

func getNewBan(dbCTX *sqlite.Context) ([]types.BanApplication, error) {

	var bas []types.BanApplication

	//select the news bans
	banRecords := dbCTX.Db.
		Order("updated_at desc").
		/*Get non expired (until) bans*/
		Where(`strftime("%s", until) >= strftime("%s", "now")`).
		/*Only get one ban per unique ip_text*/
		Group("ip_text").
		Find(&bas)
	if banRecords.Error != nil {
		return nil, fmt.Errorf("failed when selection bans : %v", banRecords.Error)
	}

	return bas, nil

}

func getLastBan(dbCTX *sqlite.Context, lastTS time.Time) ([]types.BanApplication, error) {

	var bas []types.BanApplication

	//select the news bans
	banRecords := dbCTX.Db.
		Order("updated_at desc").
		/*Get non expired (until) bans*/
		Where(`strftime("%s", until) >= strftime("%s", "now")`).
		/*That were added since last tick*/
		Where(`strftime("%s", updated_at) >= strftime("%s", ?)`, lastTS).
		/*Only get one ban per unique ip_text*/
		Group("ip_text").
		Find(&bas) /*.Count(&count)*/
	if banRecords.Error != nil {
		return nil, fmt.Errorf("failed when selection bans : %v", banRecords.Error)
	}

	return bas, nil

}

func getDeletedBan(dbCTX *sqlite.Context, lastTS time.Time) ([]types.BanApplication, error) {
	var bas []types.BanApplication

	deletedRecords := dbCTX.Db.
		/*ignore the soft delete*/
		Unscoped().
		Order("updated_at desc").
		/*ban that were deleted since lastTS or bans that expired since lastTS*/
		Where(`strftime("%s", deleted_at) >= strftime("%s", ?) OR 
		   (strftime("%s", until) >= strftime("%s", ?) AND strftime("%s", until) <= strftime("%s", "now"))`,
			lastTS.Add(1*time.Second), lastTS.Add(1*time.Second)).
		/*Only get one ban per unique ip_text*/
		Group("ip_text").
		Find(&bas) /*.Count(&count)*/

	if deletedRecords.Error != nil {
		return nil, fmt.Errorf("failed when selection deleted bans : %v", deletedRecords.Error)
	}

	return bas, nil
}
