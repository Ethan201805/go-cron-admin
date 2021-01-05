package router

import (
	"github.com/gin-gonic/gin"
	"go-cron/master/api"
)

func InitJobRouter(r *gin.Engine){

	j := r.Group("/job")
	{
		j.POST("",api.SaveJob)
		//j.PUT("",api.SaveJob)
		j.DELETE("/:name",api.DeleteJob)
	}

	l := r.Group("")
	{
		l.GET("/job/list",api.JobList)
		l.DELETE("/jobKill/:name",api.KillJob)
	}
}
