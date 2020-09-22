import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardActions from '@material-ui/core/CardActions';
import TextField from '@material-ui/core/TextField';
import CircularProgress from '@material-ui/core/CircularProgress';
import { Project, ProjectAPI } from "../../../services/project";

const useStyles = makeStyles((theme) =>({
    root: {
        minHeight: 400,
        minWidth: 400
    },
    title: {
        top: 1,
    },
    actionButtons: {
        top: '145px',
        position: 'relative'
    },
    progress: {
        display: 'flex',
        '& > * + *': {
          marginLeft: theme.spacing(2),
        },
    },
  }),
);

interface IProps {
    Callback(project: Project): void;
    ProjectAPI: ProjectAPI;
}


export function ProjectModal(props: IProps) {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [showProgress, setShowProgress] = useState(false);
    const onSubmit = function (event: React.FormEvent) {
        event.preventDefault();
    }

    const onSave = async () => {
        const proj = {
            metadata: {
                name,
                description,
                id: '',
            },
            tasks: [],
        }
        const project = await props.ProjectAPI.create(proj)
        setShowProgress(true)
        setTimeout(() => {
            props.Callback(project)
        }, 1000)
    }

    const classes = useStyles();
    if (showProgress) {
        return( 
            <div className={classes.progress}>
                <CircularProgress></CircularProgress>
            </div>
        )
    }
    return (
        <form noValidate autoComplete="off" onSubmit={onSubmit}>
            <Card className={classes.root}>
                <TextField value={name} onChange={(ev: any) => setName(ev.target.value)} fullWidth label="Name" variant="filled" />
                <CardContent>
                    <TextField
                        label="Description"
                        multiline
                        value={description}
                        onChange={(ev: any) => setDescription(ev.target.value)}
                        rows={4}
                        fullWidth
                        defaultValue="Write some task description"
                        variant="filled"
                    />
                </CardContent>
                <CardActions className={classes.actionButtons}>
                    <Button type="submit" size="small" color="primary" onClick={onSave}>
                        Save
                    </Button>
                </CardActions>
            </Card>
        </form>
    )
}