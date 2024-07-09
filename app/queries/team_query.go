package queries

import (
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/jmoiron/sqlx"
)

type TeamQueries struct {
	*sqlx.DB
}

func (q *TeamQueries) CreateNewTeam(team *models.Team) error {
	query := `
	INSERT INTO teams
	VALUES ($1, $2, $3, $4, $5)
	`

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
	query := `
	INSERT INTO team_members
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := q.Exec(
		query,
		teamMember.MemberID, teamMember.TeamID, teamMember.Role, teamMember.CreatedAt, teamMember.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *TeamQueries) CreateNewTeamInvite(teamMember *models.TeamInvite) error {
	query := `
	INSERT INTO team_invites
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := q.Exec(
		query,
		teamMember.InviteID, teamMember.TeamID, teamMember.UserID, teamMember.Role, teamMember.ExpiryDate, teamMember.CreatedAt, teamMember.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *TeamQueries) GetTeamInvite(userID, inviteID string) (models.TeamInvite, error) {
	var teamInvite models.TeamInvite
	query := `
	SELECT *
	FROM team_invites
	WHERE invite_id = $1 
		AND user_id = $2
	`

	err := q.Get(&teamInvite, query, inviteID, userID)
	if err != nil {
		return teamInvite, err
	}
	return teamInvite, nil
}

func (q *TeamQueries) GetTeamInviteByUserIDAndTeamID(teamID, userID string) (models.TeamInvite, error) {
	var teamInvite models.TeamInvite
	query := `
	SELECT *
	FROM team_invites
	WHERE team_id = $1 
		AND user_id = $2
	`

	err := q.Get(&teamInvite, query, teamID, userID)
	if err != nil {
		return teamInvite, err
	}
	return teamInvite, nil
}

func (q *TeamQueries) DeleteTeamInviteByID(inviteID string) error {
	query := `
	DELETE FROM team_invites
	WHERE invite_id = $1
	`

	_, err := q.Exec(query, inviteID)
	if err != nil {
		return err
	}
	return nil
}

func (q *TeamQueries) GetTeamMember(memberID, teamID string) (models.TeamMember, error) {
	var teamMember models.TeamMember
	query := `
	SELECT member_id,
		team_id,
		role,
		created_at,
		updated_at
	FROM team_members
	WHERE member_id = $1
		AND team_id = $2
	`

	err := q.DB.Get(&teamMember, query, memberID, teamID)
	if err != nil {
		return teamMember, err
	}
	return teamMember, nil
}

func (q *TeamQueries) GetAllTeamMembersWithDetails(teamID string) ([]models.TeamMemberWithDetails, error) {
    var teamMembers []models.TeamMemberWithDetails
    query := `
    SELECT 
        tm.member_id, 
        tm.team_id, 
        tm.role, 
        u.name, 
        u.email,
        u.avatar,
		tm.created_at,
		tm.updated_at
	FROM team_members tm
    JOIN users u ON tm.member_id = u.user_id
    WHERE tm.team_id = $1
    `

    err := q.DB.Select(&teamMembers, query, teamID)
    if err != nil {
        return teamMembers, err
    }
    return teamMembers, nil
}

func (q *TeamQueries) GetUserAllTeams(userID string) ([]models.UserTeams, error) {
    var userTeams []models.UserTeams
    query := `
    SELECT 
        tm.team_id, 
        tm.member_id, 
        t.name, 
        t.description, 
        tm.role, 
		tm.created_at,
		tm.updated_at
	FROM team_members tm
    JOIN teams t ON t.team_id = tm.team_id
    WHERE tm.member_id = $1
    `

    err := q.DB.Select(&userTeams, query, userID)
    if err != nil {
        return userTeams, err
    }
    return userTeams, nil
}
