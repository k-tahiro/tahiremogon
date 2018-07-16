package model

type (
    Command struct {
        ID   string `db:"id"`
        Name string  `db:"name"`
        Signal string  `db:"signal"`
    }
	
	Response struct {
		Success bool `json:"success"`
	}

	Request struct {
		ID string `json:"id"`
		Name string `json:"name"`
	}
)
