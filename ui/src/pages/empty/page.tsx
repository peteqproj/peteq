import React from 'react'

interface IProps {
    data: string;
}

export function EmptyPage(props: IProps) {
    return (
        <div>Empty page for: {props.data}</div>
    )
}