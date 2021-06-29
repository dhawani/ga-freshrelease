package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type IssueResponse struct {
	Issue struct {
		ID           int         `json:"id"`
		Key          string      `json:"key"`
		Title        string      `json:"title"`
		Name         string      `json:"name"`
		IssueTypeID  int         `json:"issue_type_id"`
		StatusID     int         `json:"status_id"`
		Resolved             bool        `json:"resolved"`
		Deleted              bool        `json:"deleted"`
	} `json:"issue"`
	IssueTypes []struct {
		Label       string `json:"label"`
		Description string `json:"description"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		TypeIcon    string `json:"type_icon"`
		Default     bool   `json:"default"`
		Level       string `json:"level"`
		Editable    bool   `json:"editable"`
		Removable   bool   `json:"removable"`
	} `json:"issue_types"`
	Statuses []struct {
		Label            string `json:"label"`
		Description      string `json:"description"`
		ID               int    `json:"id"`
		Name             string `json:"name"`
		StatusCategoryID int    `json:"status_category_id"`
		Estimatable      bool   `json:"estimatable"`
		Default          bool   `json:"default"`
		Removable        bool   `json:"removable"`
	} `json:"statuses"`
}

func GetIssue(baseUrl string, token string, projectKey string, issueKey string) (*IssueResponse, error) {
	url := fmt.Sprintf("%s/%s/issues/%s", baseUrl, projectKey, issueKey)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Add("Accept", "application/json")

	if res, err := http.DefaultClient.Do(req); err == nil {
		defer res.Body.Close()
		switch res.StatusCode {
		case http.StatusOK:
			iR := IssueResponse{}
			err = json.NewDecoder(res.Body).Decode(&iR)
			if err == nil {
				return &iR, nil
			} else {
				return nil, err
			}
		case http.StatusUnauthorized:
			return nil, errors.New("unauthorized")
		case http.StatusNotFound:
			return nil, errors.New(fmt.Sprintf("%s not found", issueKey))
		default:
			return nil, errors.New(fmt.Sprintf("unknown status code: %d", res.StatusCode))
		}
	} else {
		return nil, err
	}
	return nil, nil
}