package main

import (
	"errors"
	"io/ioutil"
	"path"
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
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			// return "/avatars/" + useridStr + ".jpg", nil
			files, err := ioutil.ReadDir("avatars")
			if err != nil {
				return "", ErrNoAvatarURL
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if match, _ := path.Match(useridStr+"*", file.Name()); match {
					return "/avatars/" + file.Name(), nil
				}
			}
		}
	}
	return "", ErrNoAvatarURL

}
