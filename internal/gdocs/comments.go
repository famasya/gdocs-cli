package gdocs

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Comment represents a simplified Google Docs comment.
type Comment struct {
	Author      string
	Content     string
	QuotedText  string
	CreatedTime string
	Resolved    bool
	Replies     []Reply
}

// Reply represents a reply to a comment.
type Reply struct {
	Author      string
	Content     string
	CreatedTime string
}

// FetchComments retrieves all comments for a document using the Drive API.
func FetchComments(ctx context.Context, httpClient *http.Client, docID string) ([]Comment, error) {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive service: %w", err)
	}

	var comments []Comment
	pageToken := ""
	for {
		call := srv.Comments.List(docID).Fields("comments(author(displayName),content,quotedFileContent,createdTime,resolved,replies(author(displayName),content,createdTime)),nextPageToken").PageSize(100)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		resp, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve comments: %w", err)
		}

		for _, c := range resp.Comments {
			if c.Deleted {
				continue
			}
			comment := Comment{
				Content:     c.Content,
				CreatedTime: c.CreatedTime,
				Resolved:    c.Resolved,
			}
			if c.Author != nil {
				comment.Author = c.Author.DisplayName
			}
			if c.QuotedFileContent != nil {
				comment.QuotedText = c.QuotedFileContent.Value
			}
			for _, r := range c.Replies {
				if r.Deleted {
					continue
				}
				reply := Reply{
					Content:     r.Content,
					CreatedTime: r.CreatedTime,
				}
				if r.Author != nil {
					reply.Author = r.Author.DisplayName
				}
				comment.Replies = append(comment.Replies, reply)
			}
			comments = append(comments, comment)
		}

		pageToken = resp.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return comments, nil
}
