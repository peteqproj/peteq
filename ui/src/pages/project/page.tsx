import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Chip from '@material-ui/core/Chip';
import Typography from '@material-ui/core/Typography';
import { RouteComponentProps } from "react-router-dom";
import DoneIcon from '@material-ui/icons/Done';
import UndoIcon from '@material-ui/icons/Undo';
import { Table, RowMenuItem } from './../../components/table';
import { Fab } from './../../components/fab';
import { Task as TaskComponent } from './../../components/task/task';
import { ProjectAPI } from './../../services/project'
import { Task, TaskAPI } from './../../services/tasks'
import { List, ListAPI } from './../../services/list'
import { ProjectViewAPI, ProjectView } from './../../services/view.project'
import { Progress } from './progress';

const useStyles = makeStyles({
  root: {
    display: 'flex',
  },
  content: {
    width: '100%'
  },
  media: {
    height: 320,
    width: 320
  },
});


interface IProps extends RouteComponentProps {
  id?: string;
  ProjectAPI: ProjectAPI;
  ProjectViewAPI: ProjectViewAPI;
  TaskAPI: TaskAPI;
  ListAPI: ListAPI;
}

export function ProjectPage(props: IProps) {
  const classes = useStyles();
  const projectId = (props.match.params as any)['id'];
  const [state, setState] = useState({
    metadata: {
      name: '',
      id: '',
      description: '',
      color: '',
      imageUrl: '',
    }, tasks: []
  } as ProjectView);
  useEffect(() => {
    const fetch = async () => {
      const prj = await props.ProjectViewAPI.get(projectId)
      setState(prj)
    }
    fetch();

  }, [props.ProjectAPI, projectId]);
  return (
    <Card className={classes.root} >
      <CardContent className={classes.content}>
        <CardHeader component={() => (<Chip style={{ width: '70px', backgroundColor: state.metadata.color }} />)} />
        <Typography gutterBottom variant="h3" component="h2">
          {state.metadata.name}
        </Typography>
        <Typography variant="body1" color="textSecondary" component="p">
          {state.metadata.description} {calculateCompleted(state.tasks)} / {state.tasks.length}
        </Typography>
        <Typography variant="body2" color="textSecondary" component="p">
          <Progress value={calculateProgress(state.tasks)} />
          <Table
            columns={[{ title: 'Name' }, { title: 'Description' }, { title: 'Status' }]}
            data={
              state.tasks.map(t => ({
                data: [
                  { content: t.metadata.name },
                  { content: t.metadata.description },
                  { content: t.status.completed ? 'Completed' : '' },
                ],
                menu: [
                  makeTaskCompletionButton(t, async (index: number, action: string) => {
                    if (action === 'Reopen') {
                      await props.TaskAPI.reopen(t.metadata.id);
                    }

                    if (action === 'Complete') {
                      await props.TaskAPI.complete(t.metadata.id);
                    }
                    const view = await props.ProjectViewAPI.get(projectId);
                    setState(view);
                  })
                ],
              }))
            }
          />
        </Typography>
      </CardContent>
      {state.metadata.imageUrl !== '' && <CardMedia
        className={classes.media}
        image={state.metadata.imageUrl}
        title="Contemplative Reptile"
      />}
      <Fab
        modal={
          <TaskComponent
            onChange={async () => {
              const view = await props.ProjectViewAPI.get(projectId);
              setState(view);
            }}
            defaultProject={ { name: state.metadata.name, id: state.metadata.id } }
            projects={[
              { name: state.metadata.name, id: state.metadata.id }
            ]} 
            new={true}
            task={{ metadata: { id: '', name: '', description: '' }, spec: {}, status: {completed: false}}} 
            TaskAPI={props.TaskAPI}
            ListAPI={props.ListAPI}
            ProjectAPI={props.ProjectAPI}
          />}
      />
    </Card>
  );
}

function calculateProgress(tasks: Task[]): number {
  const completed = calculateCompleted(tasks);
  const total = tasks.length;
  return (completed / total) * 100
}

function calculateCompleted(tasks: Task[]): number {
  return tasks.reduce((acc, curr) => {
    if (curr.status.completed) {
      return acc + 1
    }
    return acc;
  }, 0)
}

function makeTaskCompletionButton(task: Task, onClick: (index: number, action: string) => void): RowMenuItem {
  const icon = task.status.completed ? UndoIcon : DoneIcon;
  const action = !task.status.completed ? 'Complete' : 'Reopen';
  return {
    title: action,
    icon: icon,
    onClick: (index: number) => onClick(index, action),
  };
}  
