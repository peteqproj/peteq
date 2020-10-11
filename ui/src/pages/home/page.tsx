import React, { useEffect, useState } from 'react'
import Bluebird from 'bluebird';
import { get, cloneDeep } from 'lodash';
import { DragDropContext, Droppable, Draggable, DropResult } from 'react-beautiful-dnd';
import { makeStyles, Theme } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Container from '@material-ui/core/Container';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import IconButton from '@material-ui/core/IconButton';
import AddIcon from '@material-ui/icons/Add';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import Dialog from '@material-ui/core/Dialog';
import Backdrop from '@material-ui/core/Backdrop';
import { HomeViewAPI, HomeViewModel, HomeViewTask } from '../../services/views/home';
import { TaskAPI, Task } from "./../../services/tasks";
import { ListAPI, List } from "./../../services/list";
import { ProjectAPI, Project } from "./../../services/project";
import { Task as TaskComponent } from './../../components/task/task';
import { FullScreenDialog } from './../../components/fullscreen-dialog';


const useStyles = makeStyles((theme: Theme) => ({
    root: {
        flexGrow: 1,
    },
    paper: {
        height: 800,
        width: 300,
        backgroundColor: '#80808047',
        overflow: 'auto'
    },
    modal: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    control: {
        padding: theme.spacing(5),
    },
    listTitle: {
        paddingTop: 10,
        paddingBottom: 10,
    },
    card: {
        display: 'flex',
        maxHeight: 200,
        marginBottom: 5,
        borderLeftStyle: 'outset',
        borderLeftWidth: '10px',
    },
    progress: {
        display: 'flex',
        '& > * + *': {
            marginLeft: theme.spacing(2),
        },
    },
    addCard: {
        position: 'relative',
        bottom: '45px',
        left: '250px'
    }
}),
);

interface IProps {
    TaskAPI: TaskAPI;
    ListAPI: ListAPI;
    ProjectAPI: ProjectAPI;
    HomeViewAPI: HomeViewAPI;
}

export function HomePage(props: IProps) {
    const classes = useStyles();
    const [showNewTask, setShowNewTask] = useState(false);
    const [newTaskListIndex, setNewTaskListIndex] = useState(-1);
    const [newTaskName, setNewTaskName] = useState("");
    const [showTaskModal, setShowTaskModal] = useState(false);
    const [taskModal, setTaskModal] = useState<Task>({} as Task);
    const [state, setState] = useState<HomeViewModel>({ lists: [] });
    const [ selectedTaskListIndex, setSelectedTaskListIndex ] = useState(-1);
    const [ selectedTaskProjectIndex, setSelectedTaskProjectIndex ] = useState(-1);
    const [ projects, updateProjects ] = useState<Project[]>([]);

    useEffect(() => {
        (async () => {
            const view = await props.HomeViewAPI.get();
            setState(view);
            const projects = await props.ProjectAPI.list();
            updateProjects(projects);
        })()
    }, [])

    const onDragEnd = (drag: DropResult) => {
        const { source, destination } = drag;
        // dropped outside the list
        if (!destination) {
            return;
        }

        let task: Task
        const lists = cloneDeep(state.lists);
        lists.forEach(list => {
            if (list.metadata.id !== source.droppableId) {
                return
            }
            task = list.tasks.splice(source.index, 1)[0]
        });
        lists.forEach(list => {
            if (list.metadata.id !== destination.droppableId) {
                return
            }
            list.tasks.splice(destination.index, 0, task)
        })
        const fetchAndUpdate = async () => {
            await props.ListAPI.moveTasks(source.droppableId, destination.droppableId, [task.metadata.id])
            setState(() => ({ lists }));
        };
        fetchAndUpdate()

    }

    const onShowNewTask = (index: number) => {
        setShowNewTask(true);
        setNewTaskListIndex(index);
    }

    const onAddTask = (list: string, index: number) => {
        return async (e: any) => {
            if (e.keyCode !== 13) {
                return;
            }
            const task = await props.TaskAPI.create({
                metadata: {
                    name: newTaskName,
                    description: '',
                    id: ''
                },
                spec: {},
                status: {
                    completed: false,
                }
            });
            await Bluebird.delay(2000)
            await props.ListAPI.moveTasks('', list, [task.metadata.id]);
            setShowNewTask(true);
            setNewTaskName('')
            setState((prev) => {
                const s = cloneDeep(prev);
                s.lists[index].tasks.push(task)
                return s;
            })

        }
    };

    const onTaskClick = (task: HomeViewTask, listIndex: number) => {
        setShowTaskModal(true);
        setTaskModal(task);
        setSelectedTaskListIndex(listIndex)
        const id = task.project?.metadata.id;
        if (!id) {
            return
        }
        projects.map((p, i) => {
            if (p.metadata.id !== id) {
                return p
            }
            setSelectedTaskProjectIndex(i);
            return p
        })
    }
    return (
        <Grid container className={classes.root}>
            <Grid item xs={12}>
                <Grid container justify="center" spacing={10}>
                    <DragDropContext onDragEnd={onDragEnd}>
                        {state.lists.map((list, listIndex) => (
                            <Droppable key={listIndex} droppableId={list.metadata.id}>
                                {(provided, snapshot) => (
                                    <Grid key={list.metadata.id} item ref={provided.innerRef}>
                                        <Paper className={classes.paper} elevation={3}>
                                            <Container fixed>
                                                <div className={classes.listTitle}>{list.metadata.name}</div>
                                                {(list.tasks || []).map((task, index) => (
                                                    <Draggable
                                                        index={index}
                                                        key={task.metadata.id}
                                                        draggableId={task.metadata.id}>
                                                        {(provided, snapshot) => (
                                                            <Card
                                                                onClick={() => onTaskClick(task, listIndex)}
                                                                {...provided.draggableProps}
                                                                {...provided.dragHandleProps}
                                                                ref={provided.innerRef}
                                                                className={classes.card}
                                                                style={{
                                                                    borderLeftColor: get(task, 'project.metadata.color', 'gray'),
                                                                    ...provided.draggableProps.style
                                                                }}>
                                                                <CardContent>
                                                                    <Typography variant="body2" component="p">
                                                                        {task.metadata.name}
                                                                    </Typography>
                                                                </CardContent>
                                                            </Card>
                                                        )}
                                                    </Draggable>
                                                ))}
                                                {showNewTask && newTaskListIndex === listIndex && <Card className={classes.card}>
                                                    <CardContent>
                                                        <TextField onBlur={(e: any) => {
                                                            if (e.currentTarget.contains(e.relatedTarget)) {
                                                                return
                                                            }
                                                            setShowNewTask(false)
                                                        }} autoFocus onKeyDown={onAddTask(list.metadata.id, listIndex)} onChange={(ev: any) => setNewTaskName(ev.target.value)} value={newTaskName} />
                                                    </CardContent>
                                                </Card>}
                                            </Container>
                                        </Paper>
                                        {provided.placeholder}
                                        <IconButton onClick={() => onShowNewTask(listIndex)} aria-label="add" color="primary" className={classes.addCard}>
                                            <AddIcon />
                                        </IconButton>
                                    </Grid>
                                )}
                            </Droppable>
                        ))}
                    </DragDropContext>
                </Grid>
                <Dialog
                    aria-labelledby="transition-modal-title"
                    aria-describedby="transition-modal-description"
                    className={classes.modal}
                    open={showTaskModal}
                    onClose={() => setShowTaskModal(false)}
                    closeAfterTransition
                    BackdropComponent={Backdrop}
                    BackdropProps={{
                        timeout: 500,
                    }}
                >
                        <FullScreenDialog
                            onClose={
                                () => console.log('closed')
                            }
                            component={
                            <TaskComponent
                                onChange={async () => {
                                    const view = await props.HomeViewAPI.get()
                                    setState(view)
                                }}
                                projects={projects.map(p => ({ id: p.metadata.id, name: p.metadata.name }))}
                                defaultProject={
                                    {
                                        id: get(projects, `[${selectedTaskProjectIndex}].metadata.id`),
                                        name: get(projects, `[${selectedTaskProjectIndex}].metadata.name`),
                                    }
                                }
                                defaultList={
                                    {
                                        id: get(state, `lists[${selectedTaskListIndex}].metadata.id`),
                                        name: get(state, `lists[${selectedTaskListIndex}].metadata.name`),
                                    }
                                }
                                lists={state.lists.map(l => ({ id: l.metadata.id, name: l.metadata.name }))}
                                new={false}
                                task={taskModal}
                                TaskAPI={props.TaskAPI}
                                ListAPI={props.ListAPI}
                                ProjectAPI={props.ProjectAPI}
                            />
                        }
                    />
                </Dialog>
            </Grid>
        </Grid>
    );
}

