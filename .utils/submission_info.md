# Groupie-Tracker Submission Information

## Repository Information
- **Original GitHub Repository**: https://github.com/MiddSea/Groupie-Tracker
- **Target Submission Repository**: https://01.gritlab.ax/git/bosterma/groupie-tracker

## Submission Process
1. Björn needs to push this repository to the GitLab (Gitea) instance at GritLab
2. The repository should be submitted for audit as soon as possible

## Audit Preparation
- Both Seán David Middleton and Björn Österman need to understand the code thoroughly
- Key files to focus on:
  - `main.go`: Contains all the core functionality
  - `templates/*.html`: Frontend HTML templates
  - `static/*`: CSS and JavaScript files

## Project Structure
The project follows a standard Go web application structure:
- Main application logic in `main.go`
- HTML templates in `templates/`
- Static assets in `static/`
- Utility scripts in `utils/`

## Testing
Make sure to run tests before submission:
```bash
go test -v
```

## Key Features
- Artist information display
- Search functionality
- Responsive design
- Error handling