// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package rating

import (
	"database/sql"
)

type Bgrating struct {
	ID                 int32
	CharName           string
	CharGuildName      string
	CharLevel          int32
	SoloRating         int32
	SoloWins           int32
	SoloLoses          int32
	LastSoloMatchTime  sql.NullTime
	PartyRating        int32
	PartyWins          int32
	PartyLoses         int32
	LastPartyMatchTime sql.NullTime
	BattlesCount       int32
}
