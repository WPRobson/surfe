package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Action struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"`
	UserID     int       `json:"userId"`
	TargetUser int       `json:"targetUser"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ActionCount struct {
	Count int `json:"count"`
}

type ActionProbability map[string]float64

type ReferralIndex struct {
	Index map[int]int `json:"index"`
}

// type ReferralIndex map[int]int
