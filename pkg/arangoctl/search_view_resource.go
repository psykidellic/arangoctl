package arangoctl

import log "github.com/sirupsen/logrus"

// View implementation of resource
type SearchView struct {
	config SearchViewConfig
}

func (v *SearchView) GetKind() string {
	return "view"
}

func (v *SearchView) Apply (client *Client) error {

	// Depending upon if the view exists or not, just
	// call create or update.
	exists, err := client.ViewExists(v.config.Meta.Name)
	if err != nil {
		return err
	}

	if exists {
		log.Infof("View %s exists", v.config.Meta.Name)
		err := client.UpdateView(v.config.Meta.Name, v.config.SearchViewProperties)
		if err != nil {
			return err
		}
		log.Infof("View %s updated", v.config.Meta.Name)
	} else {
		_, err := client.CreateView(v.config.Meta.Name, &v.config.SearchViewProperties)
		if err != nil {
			return err
		}
		log.Infof("View %s created", v.config.Meta.Name)
	}

	return nil
}
