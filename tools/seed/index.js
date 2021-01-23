// Register and generate tasks/projects and backlog

const Chance = require('chance');
const axios = require('axios');
const _ = require('lodash');
const Bluebird = require('bluebird');

const memory = {
    lists: [],
    projects: [],
};

async function start() {
    setTimeout(() => process.exit(1), 60000);
    const email = process.env.EMAIL || 'demo@peteq.io';
    const password = process.env.password || 'demo';
    const url = process.env.URL || 'http://localhost'
    let token = '';
    console.log(`Starting seed process, email: ${email}, password: ${password}, url: ${url}`);
    await register({url, email, password});
    console.log('Registration request approved, waiting to be completed');
    for (let index = 0; index < 10; index++) {
        try {
            console.log(`Loggin attempt #${index}`);
            token = await login({url, email, password});
            break;
        } catch(err) {
            await Bluebird.delay(2000);
        }
    }
    console.log(`Authentication token recieved: ${token}`);
    await getLists({url, token})
    const s = new Chance();
    try {
        for (let index = 0; index < s.integer({min: 3, max: 7}); index++) {
            const project = generateProject(s);
            console.log(`Generating project: ${project.name}`);
            const pid = await createProject({url, token, project });
            memory.projects.push(pid);
        }
    } catch(err) {
        console.log(`Failed to create project: ${err.message}`);
        process.exit(1);
    }

    try {
        for (let index = 0; index < s.integer({ min: 50, max: 80 }); index++ ) {
            const task = generateTask(s);
            console.log(`Generating Task: ${task.name}`);
            await createTask({url, token, task})
        }
    } catch(err) {
        console.log(`Failed to create task: ${err.message}`);
        process.exit(1);
    }
    process.exit(0);
}

async function getLists({ url, token }) {
    return axios({
        method: 'GET',
        url: `${url}/api/list`,
        headers: {
            "content-type": "application/json",
            "authorization": token,
        }, 
    }).then(res => {
        _.map(res.data, (list) => {
            memory.lists.push({ id: list.metadata.id, name: list.metadata.name });
        });
    })
}

async function register({url, email, password, attempt = 0}) {
    try {
        if (attempt !== 0) {
            console.log(`Registration attemptt #${attempt+1}`);
        }
        await axios({
            url: `${url}/c/user/register`,
            method: 'POST',
            data: {
                email,
                password
            }
        })
    } catch (err) {
        console.log(err.message);
        await Bluebird.delay(3000);
        await register({url, email, password, attempt: attempt + 1 });
    }
    
}

async function login({url, email, password}) {
    return axios({
        url: `${url}/c/user/login`,
        method: 'POST',
        data: {
            email,
            password
        }
    })
    .then(res => {
        return _.get(res, 'data.data.token');
    })
}

function generateTask(seed) {
    const name = seed.sentence({ words: seed.integer({ min: 1, max: 5 }) });
    const description = seed.sentence({ words: seed.integer({ min: 5, max: 25 }) });
    const result = { name, description };
    if (!seed.bool()) {
        return result;   
    }
    result.list = memory.lists[seed.integer({min: 0, max: 2})].id;

    if (!seed.bool()) {
        return result;
    }

    result.project = memory.projects[seed.integer({min: 0, max: memory.projects.length })];
    return result;
}

async function createTask({url, token, task}) {
    return axios({
        url: `${url}/c/task/create`,
        headers: {
            "content-type": "application/json",
            "authorization": token,
        }, 
        method: 'POST',
        data: task,
    })
}

function generateProject(seed) {
    const name = seed.sentence({ words: seed.integer({ min: 1, max: 5 }) });
    const description = seed.sentence({ words: seed.integer({ min: 5, max: 25 }) });
    const color = seed.color();
    const imageUrl = "https://picsum.photos/200/300";
    const result = { name, description, color, imageUrl };
    return result;
}

function createProject({url, token, project}) {
     return axios({
        url: `${url}/c/project/create`,
        headers: {
            "content-type": "application/json",
            "authorization": token,
        }, 
        method: 'POST',
        data: project,
    }).then(res => {
        return res.data.id;
    })
   
}

start()

