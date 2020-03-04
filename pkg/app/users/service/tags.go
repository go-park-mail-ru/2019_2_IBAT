package users

func (h *UserService) GetTags() (map[string][]string, error) {
	tagMap := map[string][]string{}

	tags, err := h.Storage.GetTags()
	if err != nil {
		return tagMap, err
	}

	for _, item := range tags {
		if tagMap[item.ParentTag] == nil {
			tagMap[item.ParentTag] = []string{}
		}
		tagMap[item.ParentTag] = append(tagMap[item.ParentTag], item.ChildTag)
	}

	return tagMap, nil
}
