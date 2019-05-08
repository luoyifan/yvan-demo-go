define(__FILE__, function (context) {

    return [
        {
            name: 'db_ent',
            xtype: 'db',
            driver: "com.mysql.jdbc.Driver",
            url: "jdbc:mysql://47.99.62.170:3306/ent?autoReconnect=true&useUnicode=true&characterEncoding=utf-8&zeroDateTimeBehavior=convertToNull&useSSL=false",
            user: "ent",
            password: "test2018"
        },
        {
            name: 'db_wms1',
            xtype: 'db',
            driver: "com.mysql.jdbc.Driver",
            url: "jdbc:mysql://47.99.62.170:3306/wms1?autoReconnect=true&useUnicode=true&characterEncoding=utf-8&zeroDateTimeBehavior=convertToNull&useSSL=false",
            user: "wms1",
            password: "test1"
        },
    ];
});