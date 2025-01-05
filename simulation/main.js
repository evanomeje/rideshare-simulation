const wait = ms => new Promise(resolve => setTimeout(resolve, ms));

const main = async () => {
    while(true) {
        console.log(`Hello ${Math.floor(Math.random() * 100)}`);
        await wait(200);
    }
};

main();