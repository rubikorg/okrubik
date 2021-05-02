package commands

import (
	"fmt"

	"github.com/rubikorg/dice"
	"github.com/rubikorg/rubik/pkg"
	"github.com/spf13/cobra"
)

var newModel string

// TODO: add this feature when you are able to parse the model struct
// var fields []string

func initDiceCommand() *cobra.Command {
	storeCommand := &cobra.Command{
		Use:   "dice",
		Short: "Interlinks with dice ORM function generator for rubik (currently only mongo)",
		Run: func(cmd *cobra.Command, args []string) {
			// err := models.Init()
			// if err != nil {
			// 	pkg.ErrorMsg(err.Error())
			// 	return
			// }

			// ---- This is the save example for documentation
			// g := models.Game()
			// g.Name = "CSGo"
			// id := g.Save()
			// if !id.IsZero() {
			// 	fmt.Println("saved", id.String())
			// 	return
			// }

			// -- -  This is delete code for documentation
			// oid, err := primitive.ObjectIDFromHex("5fea1d833aab3badba186beb")
			// if err != nil {
			// 	pkg.ErrorMsg(err.Error())
			// 	return
			// }
			// g.FindOne(dice.Q{{"_id", oid}})
			// if !g.ID.IsZero() {
			// 	fmt.Println("I am deleting this game now haha!")
			// 	err := g.Delete()
			// 	if err != nil {
			// 		pkg.ErrorMsg(err.Error())
			// 		return
			// 	}
			// 	fmt.Println("deleted! check db")
			// }

			// newGame := models.Game()
			// newGame.Name = "CSGo"
			// id := newGame.Save()
			// if !id.IsZero() {
			// 	fmt.Println("csgo added:", id.String())
			// }

			err := execDice(cmd)
			if err != nil {
				pkg.ErrorMsg(err.Error())
			}
		},
	}

	// storeCommand.Flags().StringVarP(&genService, "gen", "g", "",
	// 	"use this flag to generate a docker files for your Rubik server")
	storeCommand.Flags().StringVarP(&newModel, "model", "m", "",
		"use this flag to create a model for your database")
	// storeCommand.Flags().StringSliceVarP(&fields, "fields", "f", []string{},
	// 	"use this flag to add fields to the model you are creating")
	return storeCommand
}

func execDice(cmd *cobra.Command) error {
	opts, err := dice.GetDiceOpts()
	if err != nil {
		return err
	}

	dice.UseOpts(opts)
	// create a model with following keys
	if newModel != "" {
		err := dice.GenerateModel(newModel)
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "insufficient arguments. nothing to do!\n")
	return nil
}
