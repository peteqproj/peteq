import React, { useEffect, useState } from 'react'
import { cloneDeep } from 'lodash';
import { DragDropContext, Droppable, Draggable, DropResult  } from 'react-beautiful-dnd';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Container from '@material-ui/core/Container';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import { TaskAPI, Task } from "./../../services/tasks";
import { ListAPI, List as ListModel } from "./../../services/list";


const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      flexGrow: 1,
    },
    paper: {
      height: 800,
      width: 300,
      backgroundColor: '#80808047'
    },
    control: {
      padding: theme.spacing(5),
    },
    listTitle: {
        paddingTop: 10,
        paddingBottom: 10,
    },
    card: {
        display: 'flex',
        height: 60,
        marginBottom: 5
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
    TaskAPI: TaskAPI;
    ListAPI: ListAPI;
}

interface List extends ListModel {
    objects: Task[]
}

interface IState {
    lists: List[]
}


export function HomePage(props: IProps) {
    const classes = useStyles();
    const [state, setState] = useState<IState>(() => {
        return {
            lists: [],
        }
      });
    
      useEffect(() => {
        (async() => {
            const lists = await props.ListAPI.list() as List[]
            const res = await props.TaskAPI.list();
            
            lists.map(list => {
                (list.tasks || []).map(id => {
                    res.map((t, i) => {
                        if (id === t.metadata.id) {
                            list.objects = (list.objects || []).concat(t);
                        };
                    })

                })
            });
            setState(() => ({ lists: lists }));
        })()
      }, [])

    const onDragEnd = (drag: DropResult) => {
        const { source, destination} = drag;
        // dropped outside the list
        if (!destination) {
            return;
        }

        let task: Task
        const lists = cloneDeep(state.lists);
        lists.forEach(list => {
            if (list.metadata.id !== source.droppableId) {
                return
            }
            task = list.objects.splice(source.index, 1)[0]
        });
        lists.forEach(list => {
            if (list.metadata.id !== destination.droppableId) {
                return
            }
            if (!list.objects) {
                list.objects = [];
            }
            list.objects.splice(destination.index, 0, task)
        })
        const fetchAndUpdate = async () => {
            await props.ListAPI.moveTasks(source.droppableId, destination.droppableId, [task.metadata.id])
            setState(() => ({ lists }));
        };        
        fetchAndUpdate()
        
    }
    return (
        <Grid container className={classes.root}>
            <Grid item xs={12}>
                <Grid container justify="center" spacing={10}>
                    <DragDropContext onDragEnd={onDragEnd}>
                        {state.lists.map((list, index) => (                
                            <Droppable key={index} droppableId={list.metadata.id}>
                                {(provided, snapshot) => (
                                    <Grid key={list.metadata.id} item ref={provided.innerRef}>
                                        <Paper className={classes.paper} elevation={3}>
                                            <Container fixed>
                                                <div className={classes.listTitle}>{list.metadata.name}</div>
                                                {(list.objects || []).map((task, index) => (
                                                    <Draggable
                                                    index={index}
                                                    key={task.metadata.id}
                                                    draggableId={task.metadata.id}>
                                                        {(provided, snapshot) => (
                                                            <Card
                                                            {...provided.draggableProps}
                                                            {...provided.dragHandleProps}
                                                            ref={provided.innerRef} 
                                                            className={classes.card}>
                                                                <CardContent>
                                                                {task.metadata.name}
                                                                </CardContent>
                                                            </Card>
                                                        )}
                                                    </Draggable>
                                                ))}
                                            </Container>
                                        </Paper>
                                    {provided.placeholder}
                                    </Grid>
                                )}
                            </Droppable>
                        ))}
                    </DragDropContext>
                </Grid>
            </Grid>
        </Grid>
    );
}

