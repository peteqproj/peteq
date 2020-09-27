import React, { useState } from 'react';
import { get } from "lodash";
import { makeStyles } from '@material-ui/core/styles';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MSelect from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';


const useStyles = makeStyles((theme) => ({
    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,
    },
}));

interface IProps {
    title: string;
    value?: string;
    items: SelectionItem[];
    onSelectionChanged(ev: SelectionChangedEvent): void;
}

interface SelectionItem {
    value: string;
    title: string;
}

export interface SelectionChangedEvent {
    destination?: {
        title: string;
        value: string;
    }
}

const UnsetSelection = {
    title: 'Not Exist',
    value: '__unset__'
}

export function Select(props: IProps) {
    const classes = useStyles();
    const [selected, setSelected] = useState(props.value || UnsetSelection.value)
    return (
        <FormControl className={classes.formControl}>
            <InputLabel >{props.title}</InputLabel>
            <MSelect
                value={selected}
                onChange={(event: React.ChangeEvent<{ name?: string; value: unknown }>, child: React.ReactNode) => {
                    const value = get(child, 'props.value', UnsetSelection.value);
                    const title = get(child, 'props.children', UnsetSelection.title);

                    if (value === UnsetSelection.value) {
                        props.onSelectionChanged({});
                        return;
                    }
                    setSelected(value);
                    props.onSelectionChanged({ destination: { value, title }});
                }}
            >
                {props.items.concat([UnsetSelection]).map((item, i) => {
                    return (
                        <MenuItem key={i} value={item.value}>{item.title}</MenuItem>
                    )
                })}
            </MSelect>
        </FormControl>
    )
}

