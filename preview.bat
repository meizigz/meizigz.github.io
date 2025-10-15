rmdir /s /q temp_hugo_content
go run converter\main.go
hugo server --contentDir temp_hugo_content -D --disableFastRender -p 4000