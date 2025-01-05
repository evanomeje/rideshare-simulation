const fs = require('fs');
const { Client } = require('pg');
const dbConfig = JSON.parse(fs.readFileSync('../dbconfig.json', 'utf8'));

const wait = ms => new Promise(resolve => setTimeout(resolve, ms));

const paths = {
    first: [
        [8, 17],
        [8, 16],
        [8, 15],
        [8, 14],
        [7, 14],
        [6, 14],
        [6, 13],
        [6, 12],
        [6, 11],
        [6, 10],
        [6, 9],
        [6, 8],
        [6, 7],
        [6, 6],
        [7, 6],
        [8, 6],
        [9, 6],
        [10, 6],
        [11, 6],
        [12, 6],
        [12, 7],
        [12, 8],
        [12, 9],
        [12, 10],
        [12, 11],
        [12, 12],
        [12, 13],
        [12, 14],
        [13, 14],
        [14, 14],
        [15, 14],
        [16, 14],
        [16, 13],
    ],
    second: [
        [16, 13],
        [16, 12],
        [16, 11],
        [16, 10],
        [15, 10],
        [14, 10],
        [13, 10],
        [12, 10],
        [11, 10],
        [10, 10],
        [9, 10],
        [8, 10],
        [8, 11],
        [8, 12],
        [8, 13],
        [8, 14],
        [8, 15],
        [8, 16],
        [8, 17],
    ]
};

const main = async () => {
    await wait(5000);

    const { host, port, user, password, dbname } = dbConfig;
    const client = new Client({ host, port, user, password, database: dbname });

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