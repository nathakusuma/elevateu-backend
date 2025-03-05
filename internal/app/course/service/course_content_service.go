package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type courseContentService struct {
	contentRepo contract.ICourseContentRepository
	courseRepo  contract.ICourseRepository
	fileUtil    fileutil.IFileUtil
	uuid        uuidpkg.IUUID
}

func NewCourseContentService(
	contentRepo contract.ICourseContentRepository,
	courseRepo contract.ICourseRepository,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.ICourseContentService {
	return &courseContentService{
		contentRepo: contentRepo,
		courseRepo:  courseRepo,
		fileUtil:    fileUtil,
		uuid:        uuid,
	}
}

func (s *courseContentService) CreateVideo(ctx context.Context, courseID uuid.UUID,
	req dto.CreateCourseVideoRequest) (dto.CreateCourseVideoResponse, error) {
	videoID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseContentService][CreateVideo] Failed to generate video ID")
		return dto.CreateCourseVideoResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	video := &entity.CourseVideo{
		ID:          videoID,
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		IsFree:      req.IsFree,
		Order:       req.Order,
	}

	err = s.contentRepo.CreateVideo(ctx, video)
	if err != nil {
		if strings.Contains(err.Error(), "course not found") {
			return dto.CreateCourseVideoResponse{}, errorpkg.ErrValidation.Build().WithDetail("Course not found")
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"video": video,
		}, "[CourseContentService][CreateVideo] Failed to create video")
		return dto.CreateCourseVideoResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// Get signed URLs for uploading video and thumbnail
	videoUploadURL, err := s.fileUtil.GetUploadSignedURL(
		fmt.Sprintf("course_videos/video/%s", videoID), "video/mp4")
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": videoID,
		}, "[CourseContentService][CreateVideo] Failed to get video upload URL")
		return dto.CreateCourseVideoResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	thumbnailUploadURL, err := s.fileUtil.GetUploadSignedURL(
		fmt.Sprintf("course_videos/thumbnail/%s", videoID), "application/pdf")
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": videoID,
		}, "[CourseContentService][CreateVideo] Failed to get thumbnail upload URL")
		return dto.CreateCourseVideoResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	response := dto.CreateCourseVideoResponse{
		CourseContent: &dto.CourseContentResponse{
			Type: "video",
			ID:   videoID,
		},
		VideoUploadURL:     videoUploadURL,
		ThumbnailUploadURL: thumbnailUploadURL,
	}

	return response, nil
}

func (s *courseContentService) UpdateVideo(ctx context.Context, id uuid.UUID, req dto.UpdateCourseVideoRequest) error {
	updates := dto.CourseVideoUpdate{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		IsFree:      req.IsFree,
		Order:       req.Order,
	}

	err := s.contentRepo.UpdateVideo(ctx, id, updates)
	if err != nil {
		if err.Error() == "video not found" {
			return errorpkg.ErrNotFound
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": id,
			"updates":  updates,
		}, "[CourseContentService][UpdateVideo] Failed to update video")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return nil
}

func (s *courseContentService) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	err := s.contentRepo.DeleteVideo(ctx, id)
	if err != nil {
		if err.Error() == "video not found" {
			return errorpkg.ErrNotFound
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": id,
		}, "[CourseContentService][DeleteVideo] Failed to delete video")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// Delete files from storage
	err = s.fileUtil.Delete(ctx, fmt.Sprintf("course_videos/video/%s", id))
	if err != nil {
		log.Error(map[string]interface{}{
			"error":    err,
			"video_id": id,
		}, "[CourseContentService][DeleteVideo] Failed to delete video file")
		// Continue execution, don't return error to client as the database deletion succeeded
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("course_videos/thumbnail/%s", id))
	if err != nil {
		log.Error(map[string]interface{}{
			"error":    err,
			"video_id": id,
		}, "[CourseContentService][DeleteVideo] Failed to delete thumbnail file")
		// Continue execution, don't return error to client as the database deletion succeeded
	}

	return nil
}

func (s *courseContentService) GetVideoUploadURLs(ctx context.Context, id uuid.UUID) (string, string, error) {
	// Check if video exists
	_, err := s.contentRepo.GetVideoByID(ctx, id)
	if err != nil {
		if err.Error() == "video not found" {
			return "", "", errorpkg.ErrNotFound
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": id,
		}, "[CourseContentService][GetVideoUploadURLs] Failed to get video")
		return "", "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	videoURL, err := s.fileUtil.GetUploadSignedURL(
		fmt.Sprintf("course_videos/video/%s", id), "video/mp4")
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": id,
		}, "[CourseContentService][GetVideoUploadURLs] Failed to get video upload URL")
		return "", "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	thumbnailURL, err := s.fileUtil.GetUploadSignedURL(
		fmt.Sprintf("course_videos/thumbnail/%s", id), "image/jpeg")
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"video_id": id,
		}, "[CourseContentService][GetVideoUploadURLs] Failed to get thumbnail upload URL")
		return "", "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return videoURL, thumbnailURL, nil
}

func (s *courseContentService) CreateMaterial(ctx context.Context, courseID uuid.UUID,
	req dto.CreateCourseMaterialRequest) (dto.CreateCourseMaterialResponse, error) {
	materialID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseContentService][CreateMaterial] Failed to generate material ID")
		return dto.CreateCourseMaterialResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	material := &entity.CourseMaterial{
		ID:       materialID,
		CourseID: courseID,
		Title:    req.Title,
		Subtitle: req.Subtitle,
		IsFree:   req.IsFree,
		Order:    req.Order,
	}

	err = s.contentRepo.CreateMaterial(ctx, material)
	if err != nil {
		if strings.Contains(err.Error(), "course not found") {
			return dto.CreateCourseMaterialResponse{}, errorpkg.ErrValidation.Build().WithDetail("Course not found")
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"material": material,
		}, "[CourseContentService][CreateMaterial] Failed to create material")
		return dto.CreateCourseMaterialResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// Get signed URL for uploading material
	materialUploadURL, err := s.fileUtil.GetUploadSignedURL(
		fmt.Sprintf("course_materials/material/%s", materialID), "application/pdf")
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"material_id": materialID,
		}, "[CourseContentService][CreateMaterial] Failed to get material upload URL")
		return dto.CreateCourseMaterialResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	response := dto.CreateCourseMaterialResponse{
		CourseContent: &dto.CourseContentResponse{
			Type: "material",
			ID:   materialID,
		},
		MaterialUploadURL: materialUploadURL,
	}

	return response, nil
}

func (s *courseContentService) UpdateCourseMaterial(ctx context.Context, id uuid.UUID,
	req dto.UpdateCourseMaterialRequest) error {
	updates := dto.CourseMaterialUpdate{
		Title:    req.Title,
		Subtitle: req.Subtitle,
		IsFree:   req.IsFree,
		Order:    req.Order,
	}

	err := s.contentRepo.UpdateMaterial(ctx, id, updates)
	if err != nil {
		if err.Error() == "material not found" {
			return errorpkg.ErrNotFound
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"material_id": id,
			"updates":     updates,
		}, "[CourseContentService][UpdateCourseMaterial] Failed to update material")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return nil
}

func (s *courseContentService) DeleteCourseMaterial(ctx context.Context, id uuid.UUID) error {
	err := s.contentRepo.DeleteMaterial(ctx, id)
	if err != nil {
		if err.Error() == "material not found" {
			return errorpkg.ErrNotFound
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"material_id": id,
		}, "[CourseContentService][DeleteCourseMaterial] Failed to delete material")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// Delete file from storage
	err = s.fileUtil.Delete(ctx, fmt.Sprintf("course_materials/material/%s", id))
	if err != nil {
		log.Error(map[string]interface{}{
			"error":       err,
			"material_id": id,
		}, "[CourseContentService][DeleteCourseMaterial] Failed to delete material file")
		// Continue execution, don't return error to client as the database deletion succeeded
	}

	return nil
}

func (s *courseContentService) GetMaterialUploadURL(ctx context.Context, id uuid.UUID) (string, error) {
	// Check if material exists
	_, err := s.contentRepo.GetMaterialByID(ctx, id)
	if err != nil {
		if err.Error() == "material not found" {
			return "", errorpkg.ErrNotFound
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"material_id": id,
		}, "[CourseContentService][GetMaterialUploadURL] Failed to get material")
		return "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	url, err := s.fileUtil.GetUploadSignedURL(
		fmt.Sprintf("course_materials/material/%s", id), "application/pdf")
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"material_id": id,
		}, "[CourseContentService][GetMaterialUploadURL] Failed to get material upload URL")
		return "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return url, nil
}

func (s *courseContentService) GetCourseContents(ctx context.Context,
	courseID uuid.UUID) ([]*dto.CourseContentResponse, error) {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	isSubscribedBoost, ok2 := ctx.Value(ctxkey.IsSubscribedBoost).(bool)
	if !ok || !ok2 {
		return nil, errorpkg.ErrInvalidBearerToken
	}

	isEnrolled := true
	_, err := s.courseRepo.GetEnrollment(ctx, courseID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "enrollment not found") {
			isEnrolled = false
			goto pass
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"course_id": courseID,
			"user_id":   userID,
		}, "[CourseContentService][GetCourseContents] Failed to get enrollment")
		return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

pass:
	isRestricted := !(isEnrolled && isSubscribedBoost)

	// (both are already sorted by order)
	videos, materials, err := s.contentRepo.GetCourseContents(ctx, courseID)
	if err != nil {
		if strings.Contains(err.Error(), "course not found") {
			return nil, errorpkg.ErrValidation.Build().WithDetail("Course not found")
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"course_id": courseID,
		}, "[CourseContentService][GetCourseContents] Failed to get course contents")
		return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	totalContents := len(videos) + len(materials)
	responses := make([]*dto.CourseContentResponse, totalContents)

	videoIdx, materialIdx, resultIdx := 0, 0, 0

	for videoIdx < len(videos) && materialIdx < len(materials) {
		// Compare orders and process the one with smaller order first
		if videos[videoIdx].Order <= materials[materialIdx].Order {
			responses[resultIdx] = &dto.CourseContentResponse{}
			err = responses[resultIdx].PopulateFromCourseVideo(videos[videoIdx], isRestricted, s.fileUtil.GetSignedURL)
			if err != nil {
				traceID := log.ErrorWithTraceID(map[string]interface{}{
					"error":    err,
					"video_id": videos[videoIdx].ID,
				}, "[CourseContentService][GetCourseContents] Failed to populate video response")
				return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
			}
			videoIdx++
		} else {
			responses[resultIdx] = &dto.CourseContentResponse{}
			err = responses[resultIdx].PopulateFromCourseMaterial(materials[materialIdx], isRestricted,
				s.fileUtil.GetSignedURL)
			if err != nil {
				traceID := log.ErrorWithTraceID(map[string]interface{}{
					"error":       err,
					"material_id": materials[materialIdx].ID,
				}, "[CourseContentService][GetCourseContents] Failed to populate material response")
				return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
			}
			materialIdx++
		}
		resultIdx++
	}

	// Process remaining
	for videoIdx < len(videos) {
		responses[resultIdx] = &dto.CourseContentResponse{}
		err = responses[resultIdx].PopulateFromCourseVideo(videos[videoIdx], isRestricted, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":    err,
				"video_id": videos[videoIdx].ID,
			}, "[CourseContentService][GetCourseContents] Failed to populate video response")
			return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
		}
		videoIdx++
		resultIdx++
	}

	for materialIdx < len(materials) {
		responses[resultIdx] = &dto.CourseContentResponse{}
		err = responses[resultIdx].PopulateFromCourseMaterial(materials[materialIdx], isRestricted,
			s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":       err,
				"material_id": materials[materialIdx].ID,
			}, "[CourseContentService][GetCourseContents] Failed to populate material response")
			return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
		}
		materialIdx++
		resultIdx++
	}

	return responses, nil
}
