# Bookings and Reservations Web App
A demo of an hotel bookings and reservations web app.

Setting Up the Source Code:
- Setup PostgreSQL Server
- Update db.config.json
- Update Makefile and/or make.bat files with correct db connection values
- Using Makefile or make.bat run the following commands: "createdb" and "migrateup"

Running the Source Code:
- Choose between 4 possible modes:
    - "ProductionMode": building the code for production
    - "DevelopmentMode": running the code on localhost with ssl disabled, 
    - "TestingMode": running all the unit tests with CSRF protection off
    - "DebuggingMode": running the code with the IDE debugger
- Change the chosen mode in the ./cmd/web/main.go file

Built Information:
- Built in Go version 1.22
- Uses Bootstrap [https://getbootstrap.com] for CSS
- Uses pgx [https://github.com/jackc/pgx] for PostgreSQL
- Uses Chi [https://github.com/go-chi/chi] for router
- Uses SCS [https://github.com/alexedwards/scs] for session management
- Uses nosurf [https://github.com/justinas/nosurf] for CSRF protection
- Uses go-simple-mail [github.com/xhit/go-simple-mail] for e-mails
- Uses testify [https://github.com/stretchr/testify] for testing
- Uses mockery [https://github.com/vektra/mockery] for mocks
- Uses  vanillajs-datepicker [https://github.com/mymth/vanillajs-datepicker] for datepicker
- Uses govalidator [https://github.com/asaskevich/govalidator] for email validation
