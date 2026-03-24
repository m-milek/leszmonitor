package handlers

import (
	"encoding/json"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
)

func CreateOrgHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.CreateOrgPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userClaims, ok := util.ExtractUserOrRespond(w, r)
	if !ok {
		return
	}

	orgCreateResponse, err := services.OrgService.CreateOrg(r.Context(), userClaims.Username, &payload)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, orgCreateResponse)
}

func DeleteOrgHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	err := services.OrgService.DeleteOrg(r.Context(), orgAuth)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "")
}

func UpdateOrgHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.UpdateOrgPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	org, err := services.OrgService.UpdateOrg(r.Context(), orgAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, org)
}

func GetOrgHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	org, err := services.OrgService.GetOrgByID(r.Context(), orgAuth)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, org)
}

func GetAllOrgsHandler(w http.ResponseWriter, r *http.Request) {
	orgs, err := services.OrgService.GetAllOrgs(r.Context())

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, orgs)
}

func AddOrgMemberHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.AddOrgMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.OrgService.AddUserToOrg(r.Context(), orgAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Member added to org successfully")
}

func RemoveOrgMemberHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.RemoveOrgMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.OrgService.RemoveUserFromOrg(r.Context(), orgAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Member removed from org successfully")
}

func ChangeOrgMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.ChangeOrgMemberRolePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.OrgService.ChangeMemberRole(r.Context(), orgAuth, payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Member role updated successfully")
}
