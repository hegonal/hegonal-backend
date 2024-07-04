package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type TeamQueries struct {
	*sqlx.DB
}

func (q *TeamQueries) CreateNewTeam(team *models.Team) error {
	query := `INSERT INTO teams VALUES ($1, $2, $3, $4, $5)`

	_, err := q.Exec(
		query,
		team.TeamID, team.Name, team.Description, team.CreatedAt, team.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *TeamQueries) CreateNewMember(teamMember *models.TeamMember) error {
	query := `INSERT INTO team_members VALUES ($1, $2, $3, $4, $5)`

	_, err := q.Exec(
		query,
		teamMember.MemberID, teamMember.TeamID, teamMember.Role, teamMember.CreatedAt, teamMember.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *TeamQueries) GetTeamMember(memberID, teamID string) (models.TeamMember, error) {
	var teamMember models.TeamMember
	query := `
		SELECT member_id, team_id, role, created_at, updated_at
		FROM team_members
		WHERE member_id = $1 AND team_id = $2
	`
	err := q.DB.Get(&teamMember, query, memberID, teamID)
	if err != nil {
		return teamMember, err
	}
	return teamMember, nil
}