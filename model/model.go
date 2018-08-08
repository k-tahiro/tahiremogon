package model

type (
    Command struct {
        ID   string `db:"id"`
        Name string  `db:"name"`
        Signal string  `db:"signal"`
    }

    CommandJSON struct {
        ID   string `json:"id"`
        Name string  `json:"name"`
        Signal string  `json:"signal"`
    }
	
	Response struct {
		Success bool `json:"success"`
	}

	Request struct {
		ID string `json:"id"`
		Name string `json:"name"`
    }
)
