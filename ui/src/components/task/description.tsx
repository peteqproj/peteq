import React, { useState, useRef } from 'react';
import { get } from 'lodash';
import CardContent from '@material-ui/core/CardContent';
import TextField from '@material-ui/core/TextField';


interface IProps {
    description: string;
    new?: boolean;
    disableAutoFocus?: boolean;
    onUpdate(description: string): void;
}

export function TaskDescription(props: IProps) {
    console.log(props.description)
    const input = useRef();
    const [editMode, setEditMode] = useState(props.new);
    const [taskDescription, setTaskDescription] = useState(props.description);
    if (editMode) {
        return (
            <CardContent
                component={(rprops) => (
                    <TextField
                        inputRef={input}
                        {...rprops}
                        onBlur={() => {
                            setEditMode(false)
                            const value = get(input, 'current.value', taskDescription);
                            setTaskDescription(value);
                            props.onUpdate(value);
                        }}
                        onKeyDown={(e: any) => {
                            if((e.ctrlKey || e.metaKey) && e.keyCode === 13){
                                setEditMode(false)
                                const value = get(input, 'current.value', taskDescription);
                                setTaskDescription(value);
                                props.onUpdate(value);
                            }
                        }}
                        multiline
                        rows={4}
                        defaultValue={taskDescription}
                        variant="outlined"
                        fullWidth
                        autoFocus={!props.disableAutoFocus}
                        placeholder="Add more details"
                />
        )} />)
    }
    return (
        <CardContent onClick={() => setEditMode(true)}>
            {taskDescription}
        </CardContent>
    )
}