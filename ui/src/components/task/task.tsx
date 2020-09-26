import React, { useState } from 'react';
import { get, cloneDeep } from 'lodash';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';
import { Task as TaskModal, TaskAPI } from './../../services/tasks';
import { List as ListModal, ListAPI } from './../../services/list';
import { Project as ProjectModal, ProjectAPI } from './../../services/project';
import { TaskTitle } from './title';
import { TaskDescription } from './description';
import { Select, SelectionChangedEvent } from './../select';

const useStyles = makeStyles((theme) => ({
    root: {
        minWidth: '400px'
    },
    content: {
        display: 'flex'
    }
}));

interface IProps {
    task: TaskModal;
    TaskAPI: TaskAPI;
    ListAPI: ListAPI;
    ProjectAPI: ProjectAPI;
    projects: { name: string, id: string }[];
    defaultProject?: {
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
    const [project, updateProject] = useState(props.defaultProject?.id || "");
    return (
        <Card className={classes.root}>
            <TaskTitle
                new={props.new}
                title={task.metadata.name}
                onUpdate={async (title: string) => {
                    const clone = cloneDeep(task)
                    clone.metadata.name = title;
                    const t = await upsert(clone, props.TaskAPI)
                    await updateTask(t)
                    if (!projectHasSet && project !== "") {
                        await props.ProjectAPI.addTasks(project, [t.metadata.id])
                        updateProjectHasSet(true);
                    }
                    props.onChange();
                }}
            />
            <Grid className={classes.content}>
                <TaskDescription
                    new={props.new}
                    disableAutoFocus={props.new}
                    description={task.metadata.description}
                    onUpdate={async (description: string) => {
                        const clone = cloneDeep(task)
                        clone.metadata.description = description;
                        const t = await upsert(clone, props.TaskAPI)
                        await updateTask(t)
                        if (!projectHasSet && project !== "") {
                            await props.ProjectAPI.addTasks(project, [t.metadata.id])
                            updateProjectHasSet(true);
                        }
                        props.onChange();
                    }}
                />
                <Select
                    onSelectionChanged={async (ev: SelectionChangedEvent) => {
                        if(task.metadata.id === "") {
                            return
                        }
                        const destination = get(ev, 'destination.value');
                        updateProject(destination);
                        updateProjectHasSet(false);
                        props.onChange();
                    }}
                    key={props.defaultProject?.name}
                    value={props.defaultProject?.id}
                    title={"Projects"}
                    items={props.projects.map(l => ({ title: l.name, value: l.id }))}
                />
            </Grid>
        </Card>
    )
}

async function upsert(task: TaskModal, api: TaskAPI): Promise<TaskModal> {
    if (task.metadata.id === "") {
        return api.create(task); 
    }
    await api.update(task);
    return task
}