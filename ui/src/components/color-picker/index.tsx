import React from 'react';
import { CirclePicker, ColorResult } from 'react-color'

interface IProps {
    onColorSet(hex: string): void;
}

export function ColorPicker(props: IProps) {
    const onChange = (ev: ColorResult) => {
        props.onColorSet(ev.hex);
    }
    return (
        <CirclePicker onChange={onChange}/>
    )
}