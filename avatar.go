package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL. ")

// Avatar reperents types capble of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if somingthing is wrong
	// ErroNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	// if url, ok := c.userData["avatar_url"]; ok {
	// 	if urlStr, ok := url.(string); ok {
	// 		return urlStr, nil
	// 	}
	// }
	// return "", ErrNoAvatarURL
	url, getURLOk := c.userData["avatar_url"]
	if !getURLOk {
		return "", ErrNoAvatarURL
	}
	urlStr, parseURLOk := url.(string)
	if !parseURLOk {
		return "", ErrNoAvatarURL
	}
	return urlStr, nil

}

type GravatarAvatar struct{}

var UseGrAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}