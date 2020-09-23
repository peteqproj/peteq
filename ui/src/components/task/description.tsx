import React, { useState, useRef } from 'react';
import { get } from 'lodash';
import CardContent from '@material-ui/core/CardContent';
import TextField from '@material-ui/core/TextField';


interface IProps {
    description: string;
    onUpdate(description: string): void;
}

export function TaskDescription(props: IProps) {
    const input = useRef();
    const [editMode, setEditMode] = useState(false);
    const [taskDescription] = useState(props.description);
    if (editMode) {
        return (
            <CardContent
                component={(rprops) => (
                    <TextField
                        inputRef={input}
                        {...rprops}
                        onBlur={() => {
                            setEditMode(false)
                            props.onUpdate(taskDescription)
                        }}
                        onKeyDown={(e: any) => {
                            if((e.ctrlKey || e.metaKey) && e.keyCode === 13){
                                setEditMode(false)
                                props.onUpdate(get(input, 'current.value', taskDescription))
                            }
                        }}
                        multiline
                        rows={4}
                        defaultValue={taskDescription}
                        variant="outlined"
                        fullWidth
                        autoFocus
                />
        )} />)
    }
    return (
        <CardContent onClick={() => setEditMode(true)}>
            {taskDescription}
        </CardContent>
    )
}