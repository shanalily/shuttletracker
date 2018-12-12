package cmd

import (
	"github.com/wtg/shuttletracker/schedule"
    "github.com/spf13/cobra"
)
// set flag to link, set default to Fall2018/Spring2019 shuttle
func init(){
	scheduleCmd.Flags().String("link", "https://rpi.box.com/shared/static/naf8gm1wjdor8tbebho5k0t28wksaygd.xlsx", "linked file")
	rootCmd.AddCommand(scheduleCmd)
}
// actual command 
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Long:  "Add schedules based on a user specified link",
	Run: func(cmd *cobra.Command, args []string) {
		link := cmd.Flag("link").Value.String()
		schedule.ReadLink(link)
	},
}
