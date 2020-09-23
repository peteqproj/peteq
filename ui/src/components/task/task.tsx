import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import { Task as TaskModal, TaskAPI } from './../../services/tasks';
import { TaskTitle } from './title';
import { TaskDescription } from './description';

const useStyles = makeStyles((theme) => ({
    root: {
      minWidth: '400px'
    },
  }));

interface IProps {
    task: TaskModal;
    TaskAPI: TaskAPI;
}


export function Task(props: IProps) {
    const classes = useStyles()
    const [task, updateTask] = useState(props.task);
    return (
        <Card className={classes.root}>
            <TaskTitle
                title={task.metadata.name}
                onUpdate={async (title: string) => {
                    await updateTask((prev) => {
                        prev.metadata.name = title;
                        return prev;
                    })
                    await props.TaskAPI.update(task)
                }}
            />
            <TaskDescription
                description={task.metadata.description}
                onUpdate={async (description: string) => {
                    await updateTask((prev) => {
                        prev.metadata.description = description;
                        return prev;
                    })
                    await props.TaskAPI.update(task)
                }}
            />
        </Card>
    )
}