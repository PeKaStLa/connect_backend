# connect_go_backend
go backend for connect app

run server with: go run main.go

## Plan:
- Web-/App where everyone has an area in every location and can chat without making a group or exchanging number or usernames...
- user can filter their area for male/female, age, verified
- users can decide if other users can see/find them
- users can see other users in their specific area and also other areas that are available at their location....
- users can subscribe to areas to still being able to chat even when they move away from them...
- users can connect/follow/friendship with other users and chat privately...
- later user can connect via near-bluetooth
- evtl. function that a few users on one spot automatically create/connect to one area? 
- or every user has an own area and when users come close together they connect somehow into one area??? for example a concert or comic con there would be no need to create one are for those locations...(but locations still make sense for hostels... or the concert location.)
- maybe user can subscribe to a location before being there?? concert??
- areas have different categories..
- the areas need to know which users are in them...
- go backend using R-tree or better R*-tree cause fast
- ...
- frontend Apps und Web in flutter cause 3-in-1
- backend with go and supabase (inclusive authentication and postgresDB)
- for infra use "just" command runner and free OKD (OpenShift cluster inclusive cicd and dev tools)
