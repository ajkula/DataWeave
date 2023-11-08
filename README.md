# README

**db_meta** is a `desktop RGBDS tool` to visualise `tables and relations graph`, `check schema integrity` (in progress), `API and SWAGGER doc generator` (todo),
`generate requests` (todo), and other yet to be designed utilities.

## About DataWeave

Simply connect to your RGBDS on the connection page, you'll be redirected to the structural graph page if everything was fine, from there go to any pages by the left nav menu.

## The project

Basically you have the `databases` package exporting the `database_manager` abstraction layer.
Its methods are attached to `(a *App) YouCallIt` to be added as Go-JS bindings.
Upon succeeding connection, the app queries all it needs into structs, and several functions are applied to transform the data to other datatypes `Marshallable` to JSON.

You can easily add the database you want by adding all the methods you can find in its own connector package, then add its support to the `database_manager` switch case and in the `connection module` selector options.

To include another page to the frontend, make the module, it needs to export:
 - `html`: contains the HTML tree, any text should have translations (todo) filled in the `lang/translations.go` package exported to the front, and it's values matched in the HTML set as `string:varName;` to be injected by the `router`, all `[varNames]` must have their matching key values in...
 - `getTranslations(pageTitle)` exports the `[varNames]` to be injected to the HTML by the router, will be done every time DOM is modified so you don't need to worry about `Tab content`.
 - `init()` method is the actual module logic, it's called when router renders the module.

You will then need to add module name into the `pagesKeys` collection for the router to know what module is called.
`Pages` contains the `page name` translation displayed.

The next steps are to add it to the router's `loadPage` switch case, and nav's `items` array :
```javascript
    {
      text: 'Graph visualizer', // Nav button name
      page: pagesKeys.graph,    // Module name to be loaded by the router
      name: pagesKeys.graph,    // temporarily used to display nav item selected
    },
```

## Project structure

Project files structure :

```javascript
db_meta/
    +- main.go
    +- app.go                       // Exports Go<->JS bindings
    +-- frontend/
    |   +-- package.json
    |   +-- index.html
    |   +-- src/
    |       +-- app.css
    |       +-- main.js             // App router, DM updates trigger vars inject
    |       +-- style.css
    |       +-- utils/
    |           +-- utils.js        // JS toolbox
    |       +-- pages/
    |           +-- connection/     // DB connection form
    |               +-- script.js   // Module logic
    |               +-- styles.css  // Module styles
    |           +-- nav/            // Nav menu using the custom router
    |               +-- script.js   
    |               +-- styles.css  
    |           +-- graph/          // Displays the Tables and relations
    |               +-- script.js   
    |               +-- styles.css  
    |           +-- integrity/      // Schema integrity checks page
    |               +-- script.js   
    |               +-- styles.css  
    |           +-- other pages soon...
    +-- databases/
    |   +-- database_connector.go   // RGBDS Interface to abstract connectors
    |   +-- database_manager.go     // Concrete implementation
    |   +-- dbstructs/
    |       +-- dbstructs.go        // Structs package
    |   +-- postgres/
    |       +-- postgres.go         // PostgreSQL connector
    |       +-- MySQL.go            // MySQL connector (todo)
    |       +-- ....go              // Other SQL connector
    |   +-- mysql/
    |       +-- mysql.go            // MySQL connector (todo)
    |   +-- sqlite/
    |       +-- sqlite.go           // SQLite connector (todo)
    |   +-- sqlserver/
    |       +-- sqlserver.go        // SQL Server connector (todo)
    +-- lang/
        +-- translations.go         // Translations package (todo)
```



## About Wails

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

To run in live development mode, run `wails dev` in the project directory.
there is also a dev server that runs on http://localhost:34115.

To build a redistributable, production mode package, use `wails build`.

