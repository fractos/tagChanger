# tagChanger
CLI tool that changes a yaml value in a github repository 

It uses the github contents api https://docs.github.com/en/rest/reference/repos#contents to download the file

and the commit it back to the repositry 

# Uses
tagChanger --file-path path --repo owner/reposiorty --branch main --commit-msg msg --value-path foo.bar --new-value new
