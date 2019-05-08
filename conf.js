define(__FILE__, function (context) {

    var ds = require('/datasource.js')(context);

    return {
        datasource: ds
    }
});