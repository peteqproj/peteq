import React, { useState } from 'react';
import { get, cloneDeep } from 'lodash';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import { Task as TaskModal, TaskAPI } from './../../services/tasks';
import { List as ListModal, ListAPI } from './../../services/list';
import { Project as ProjectModal, ProjectAPI } from './../../services/project';
import { TaskTitle } from './title';
import { TaskDescription } from './description';
import { Select, SelectionChangedEvent } from './../select';

const useStyles = makeStyles((theme) => ({
    root: {
        minWidth: '400px',
        minHeight: '100%',
        flexGrow: 1,
        paddingTop: '20px',
        backgroundColor: theme.palette.grey[100],
        margin: theme.spacing(2),
    },
    content: {
        display: 'flex'
    },
    selection: {
        width: '300px'
    },
    paper: {
        padding: '20px',
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
}));

interface IProps {
    task: TaskModal;
    TaskAPI: TaskAPI;
    ListAPI: ListAPI;
    ProjectAPI: ProjectAPI;
    projects: { name: string, id: string }[];
    lists: { name: string, id: string }[];
    defaultProject?: {
        name: string;
        id: string;
    }
    defaultList?: {
        name: string;
        id: string;
    }
    new?: boolean;
    onChange(): void;
}


export function Task(props: IProps) {
    const classes = useStyles()
    const [task, updateTask] = useState(props.task);
    const [projectHasSet, updateProjectHasSet] = useState(false);
    const [project] = useState(props.defaultProject?.id || "");
    const [listHasSet, updateListHasSet] = useState(false);
    const [selectedProject, updateSelectedProject] = useState({
        id: props.defaultProject?.id || "",
        name: props.defaultProject?.name || ""
    })
    const [selectedList, updateSelectedList] = useState({
        id: props.defaultList?.id || "",
        name: props.defaultList?.name || "",
    })

    const updateTaskList = async (task: string, list: string) => {
        if (listHasSet) {
            return; // was set previously, no changes
        }
        updateListHasSet(true);
    }
    return (
        <Card className={classes.root}>
            <Paper elevation={3} className={classes.paper}>
                <Grid container spacing={3}>
                    <Grid item xs={1}>
                        Name:
                    </Grid>
                    <Grid item xs={11}>
                        <TaskTitle
                            new={props.new}
                            title={task.metadata.name}
                            onUpdate={async (title: string) => {
                                const clone = cloneDeep(task)
                                if (!title || title === "") {
                                    return
                                }
                                clone.metadata.name = title;
                                updateTask(clone)
                                if (!projectHasSet && project !== "") {
                                    updateProjectHasSet(true);
                                }
                                props.onChange();
                            }}
                        />
                        <Divider />
                    </Grid>

                    <Grid item xs={1}>
                        Project:
                    </Grid>
                    <Grid item xs={5}>
                        <Select
                            disabled={!!selectedProject.id}
                            className={classes.selection}
                            onSelectionChanged={async (ev: SelectionChangedEvent) => {
                                const id = get(ev, 'destination.value');
                                const name = get(ev, 'destination.title');
                                updateSelectedProject({ id, name });
                                updateProjectHasSet(false);
                                props.onChange();
                            }}
                            key={selectedProject.name}
                            value={selectedProject.id}
                            items={props.projects.map(l => ({ title: l.name, value: l.id }))}
                        />
                    </Grid>

                    <Grid item xs={1}>
                        List:
                    </Grid>
                    <Grid item xs={5}>
                        <Select
                            className={classes.selection}
                            onSelectionChanged={async (ev: SelectionChangedEvent) => {
                                const id = get(ev, 'destination.value');
                                const name = get(ev, 'destination.title');
                                updateListHasSet(false);
                                updateSelectedList({ id, name })
                                await updateTaskList(task.metadata.id, id)
                                props.onChange();
                            }}
                            key={selectedList.name}
                            value={selectedList.id}
                            items={props.lists.map(l => ({ title: l.name, value: l.id }))}
                        />
                    </Grid>


                    <Grid item xs={2}>
                        Description
                    </Grid>
                    <Grid item xs={10}>
                        <TaskDescription
                            new={props.new}
                            disableAutoFocus={props.new}
                            description={task.metadata.description}
                            onUpdate={async (description: string) => {
                                const clone = cloneDeep(task)
                                const title = clone.metadata.name;
                                if (!title || title === "") {
                                    return
                                }
                                clone.metadata.description = description;
                                if (!projectHasSet && project !== "") {
                                    updateProjectHasSet(true);
                                }
                                props.onChange();
                            }}
                        />
                    </Grid>                      
                    <Grid item xs={12}>
                        <Button variant="contained" color="primary" onClick={async () => {
                                if (task.metadata.id === "") {
                                    const t = await props.TaskAPI.create({
                                        name: task.metadata.name,
                                        list: selectedList.id,
                                        project: selectedProject.id,
                                    });
                                    await updateTask(t)
                                }
                                await props.TaskAPI.update(task);
                                return task
                        }} >
                            Create
                        </Button>
                    </Grid>
                </Grid>
            </Paper>
        </Card>
    )
}