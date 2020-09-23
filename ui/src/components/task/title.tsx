import React, { useState } from 'react';
import CardHeader from '@material-ui/core/CardHeader';
import TextField from '@material-ui/core/TextField';


interface IProps {
    title: string;
    onUpdate(title: string): void;
}

export function TaskTitle(props: IProps) {
    const [editMode, setEditMode] = useState(false);
    const [taskTitle, setTaskTitle] = useState(props.title);
    if (editMode) {
        return (<CardHeader component={(rprops) => (
            <TextField 
            {...rprops}
            autoFocus
            onChange={(ev: any) => setTaskTitle(ev.target.value)}
            onBlur={() => {
                setEditMode(false)
                props.onUpdate(taskTitle)
            }}
            onKeyDown={(e: any) => {
                if(e.keyCode !== 13){ // 13 is enter
                    return;
                }
                setEditMode(false)
                props.onUpdate(taskTitle)
            }}
            value={taskTitle}
            fullWidth/>
        )} />)
    }
    return (<CardHeader title={taskTitle} onClick={() => setEditMode(true)}/>)
}