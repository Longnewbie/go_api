routes -> handler -> service -> repository -> model

config
user(model)
repo = user_repository(user)
service = user_service(repo)
handler = user_handler(service)
routes = user_routes(handler)
