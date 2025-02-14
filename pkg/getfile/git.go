package getfile

import (
	"encoding/base64"

	"github.com/xanzy/go-gitlab"
)

func GitFile(server, token, repo, file, ref string) ([]byte, error) {
	var data []byte

	clientOpts := []gitlab.ClientOptionFunc{}

	if server != "" {
		clientOpts = append(clientOpts, gitlab.WithBaseURL("https://"+server))
	}
	git, err := gitlab.NewClient(token, clientOpts...)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if ref == "" {
		log.Debug("Ref is empty, getting project settings for ", repo)
		project, _, err := git.Projects.GetProject(repo, &gitlab.GetProjectOptions{})
		if err != nil {
			return nil, err
		}
		log.Debug("Setting Ref to: ", project.DefaultBranch, " for: ", file)
		ref = project.DefaultBranch
	}

	log.Debugf("retreiving file %s from %s:%s...\n", file, repo, ref)
	gf := &gitlab.GetFileOptions{
		Ref: gitlab.String(ref),
	}
	f, _, err := git.RepositoryFiles.GetFile(repo, file, gf)
	if err != nil {
		return nil, err
	}

	// log.Debugf("File encoding: %s", f.Encoding)
	// log.Debugf("File sha256: %s", f.SHA256)
	// log.Debugf("File size: %d (bytes)", f.Size)

	if f.Encoding == "base64" {
		// log.Debug("Decoding base64...")
		data, err = base64.StdEncoding.DecodeString(f.Content)
		if err != nil {
			return nil, err
		}
	} else {
		data = []byte(f.Content)
	}

	return data, nil
}
