package model

type (
    Command struct {
        ID   int     `db:"id"`
        Name string  `db:"name"`
        // Signal string  `db:"signal"`
    }

    CommandJSON struct {
        ID   int     `json:"id"`
        Name string  `json:"name"`
        // Signal string  `json:"signal"`
	}
	
	Response struct {
		Success bool `json:"success"`
	}

	Request struct {
		Name string `json:"name"`
	}
)
