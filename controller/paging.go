package controller

// this file contains some paging related utility functions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabric8-services/admin-console/app"
	"github.com/fabric8-services/fabric8-common/log"

	errs "github.com/pkg/errors"
)

const (
	pageSizeDefault = 10
	PageSizeMax     = 10
)

func computePagingLimits(offsetParam int, limitParam int) (offset int, limit int) {
	if offsetParam < 0 {
		offset = 0
	} else {
		offset = offsetParam
	}
	limit = limitParam
	if limit <= 0 {
		limit = pageSizeDefault
	} else if limit > PageSizeMax {
		limit = PageSizeMax
	} else {
		limit = limitParam
	}
	return offset, limit
}

func setPagingLinks(links *app.PagingLinks, path string, currentCount, pageNumber, pageSize, totalCount int, additionalQuery ...string) {
	log.Info(nil, map[string]interface{}{
		"path":          path,
		"current_count": currentCount,
		"page_number":   pageNumber,
		"page_size":     pageSize,
		"total_count":   totalCount,
	}, "generating pagination links")

	format := func(additional []string) string {
		if len(additional) > 0 {
			return "&" + strings.Join(additional, "&")
		}
		return ""
	}

	// prev link
	if currentCount > 0 && pageNumber > 0 && totalCount > 0 {
		prev := fmt.Sprintf("%s?page[start]=%d&page[size]=%d%s", path, (pageNumber - 1), pageSize, format(additionalQuery))
		links.Prev = &prev
	}

	// next link
	nextStart := pageNumber + currentCount
	if currentCount > 0 && nextStart < totalCount {
		// we have a next link
		next := fmt.Sprintf("%s?page[start]=%d&page[size]=%d%s", path, (pageNumber + 1), pageSize, format(additionalQuery))
		links.Next = &next
	}

	// first link
	first := fmt.Sprintf("%s?page[start]=%d&page[size]=%d%s", path, 0, pageSize, format(additionalQuery))
	links.First = &first

	// last link
	var lastStart, lastLimit int
	if (totalCount % pageSize) == 0 {
		lastStart = (totalCount / pageSize) - 1
		lastLimit = pageSize
	} else {
		lastStart = (totalCount / pageSize)
		lastLimit = pageSize
	}
	last := fmt.Sprintf("%s?page[start]=%d&page[size]=%d%s", path, lastStart, lastLimit, format(additionalQuery))
	links.Last = &last
}

func parseInts(s *string) ([]int, error) {
	if s == nil || len(*s) == 0 {
		return []int{}, nil
	}
	split := strings.Split(*s, ",")
	result := make([]int, len(split))
	for index, value := range split {
		converted, err := strconv.Atoi(value)
		if err != nil {
			return nil, errs.WithStack(err)
		}
		result[index] = converted
	}
	return result, nil
}

func parseLimit(pageParameter *string) (s *int, l int, e error) {
	params, err := parseInts(pageParameter)
	if err != nil {
		return nil, 0, errs.WithStack(err)
	}

	if len(params) > 1 {
		return &params[0], params[1], nil
	}
	if len(params) > 0 {
		return nil, params[0], nil
	}
	return nil, 100, nil
}
