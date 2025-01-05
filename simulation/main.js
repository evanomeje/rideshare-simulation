const fs = require('fs');
const { Client } = require('pg');

const wait = ms => new Promise(resolve => setTimeout(resolve, ms));

const paths = {
    first: [
        [8,17],
        [8,16],
        [8,15],
        [9,15],
        [10,15],
        [11,15],
        [12,15],
        [13,15],
        [14,15],
        [15,15],
        [16,15],
        [16,14],
        [16,13],
        [16,12]
    ],
    second: [
        [16,12],
        [16,11],
        [16,10],
        [15,10],
        [14,10],
        [13,10],
        [12,10],
        [11,10],
        [10,10],
        [9,10],
        [8,10],
        [8,11],
        [8,12],
        [8,13],
        [8,14],
        [8,15],
        [8,16],
        [8,17]
    ]
};

const main = async () => {
    await wait(5000);

    const dbConfig = JSON.parse(fs.readFileSync('../dbconfig.json', 'utf8'));
    const { host, port, user, password, dbname } = dbConfig;
    const client = new Client({ 
        host, 
        port, 
        user, 
        password, 
        database: dbname 
    });

    client.connect((err) => {
        if (err) console.error('connection error', err.stack);
        else console.log('connected');
    });

    let path = 'first';
    let i = 0;
    while (true) {
        const [x, y] = paths[path][i];

        const res = await client.query(
            `
            INSERT INTO rides (car_id, location, path) 
            VALUES ('car1', '${x}:${y}', '${JSON.stringify(paths[path])}')
            ON CONFLICT (car_id) 
            DO UPDATE SET location = EXCLUDED.location, path = EXCLUDED.path;
            `
        );
        if (res.rowCount) console.log(`${x}:${y}`);

        if (i === paths[path].length - 1) {
            path = path === 'first' ? 'second' : 'first';
            i = 0;
            await wait(3000);
        } else {
            i++;
        }
        await wait(200);
    }
};

main().catch(console.error);
