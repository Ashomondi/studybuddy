import { useEffect, useState } from 'react'
import { api } from '../api/client'
import Card from '../components/Card'
import Loading from '../components/Loading'
import EmptyState from '../components/EmptyState'
import Badge from '../components/Badge'
import Button from '../components/Button'

const dateString = d => `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
export default function Schedule(){
 const [week,setWeek]=useState(()=>{const d=new Date();d.setDate(d.getDate()-d.getDay());return d});const [items,setItems]=useState([]);const [loading,setLoading]=useState(true);const [error,setError]=useState('')
 const end=new Date(week);end.setDate(end.getDate()+6)
 useEffect(()=>{let live=true;setLoading(true);api.schedule(dateString(week),dateString(end)).then(x=>{if(live){setItems(x.activities||[]);setError('')}}).catch(e=>live&&setError(e.message)).finally(()=>live&&setLoading(false));return()=>{live=false}},[dateString(week),dateString(end)])
 const days=Array.from({length:7},(_,i)=>{const d=new Date(week);d.setDate(d.getDate()+i);return d})
 return <div><div className="page-header"><div><h1>Study schedule</h1><p>Activities scheduled for this week.</p></div><div className="row"><Button variant="secondary" onClick={()=>{const d=new Date(week);d.setDate(d.getDate()-7);setWeek(d)}}>Previous week</Button><Button variant="secondary" onClick={()=>setWeek(()=>{const d=new Date();d.setDate(d.getDate()-d.getDay());return d})}>This week</Button><Button variant="secondary" onClick={()=>{const d=new Date(week);d.setDate(d.getDate()+7);setWeek(d)}}>Next week</Button></div></div>
 {error&&<div className="toast error" role="alert">{error}</div>}{loading?<Loading/>:items.length===0?<Card><EmptyState title="Nothing scheduled" description="Add dates to activities in a learning plan to see them here."/></Card>:<div className="grid schedule-week">{days.map(day=>{const key=dateString(day),list=items.filter(a=>a.scheduled_date===key);return <Card key={key}><h3>{day.toLocaleDateString(undefined,{weekday:'long'})}</h3><div className="muted">{day.toLocaleDateString()}</div>{list.length?list.map(a=><div className="schedule-activity" key={a.id}><b>{a.title}</b><span className="muted">{a.estimated_minutes?`${a.estimated_minutes} min`:''}</span><Badge status={a.status}/></div>):<p className="muted">No activities</p>}</Card>})}</div>}</div>
}
